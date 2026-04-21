#!/usr/bin/env python3
import json
import os
import sys
import time
from datetime import datetime, timedelta, timezone
from typing import Dict, Iterable, List, Optional

import requests


ELASTIC_URL = os.getenv("ELASTIC_URL", "http://localhost:9200").rstrip("/")
GITHUB_API_URL = os.getenv("GITHUB_API_URL", "https://api.github.com").rstrip("/")
GITHUB_OWNER = os.getenv("GITHUB_OWNER")
GITHUB_REPO = os.getenv("GITHUB_REPO")
GITHUB_TOKEN = os.getenv("GITHUB_TOKEN")
LOOKBACK_DAYS = int(os.getenv("LOOKBACK_DAYS", "180"))
REQUEST_TIMEOUT_SECONDS = int(os.getenv("REQUEST_TIMEOUT_SECONDS", "30"))

ISSUES_INDEX = "github_issues"
PULLS_INDEX = "github_pull_requests"
COMMITS_INDEX = "github_commits"


def fail(message: str) -> None:
    print(f"ERROR: {message}")
    sys.exit(1)


def parse_datetime(value: Optional[str]) -> Optional[datetime]:
    if not value:
        return None
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def duration_hours(start: Optional[str], end: Optional[str]) -> Optional[float]:
    start_dt = parse_datetime(start)
    end_dt = parse_datetime(end)
    if not start_dt or not end_dt:
        return None
    return round((end_dt - start_dt).total_seconds() / 3600, 2)


def github_headers() -> Dict[str, str]:
    headers = {
        "Accept": "application/vnd.github+json",
        "X-GitHub-Api-Version": "2022-11-28",
    }
    if GITHUB_TOKEN:
        headers["Authorization"] = f"Bearer {GITHUB_TOKEN}"
    return headers


def github_get(url: str, params: Optional[Dict] = None) -> requests.Response:
    response = requests.get(
        url,
        headers=github_headers(),
        params=params,
        timeout=REQUEST_TIMEOUT_SECONDS,
    )
    if response.status_code >= 400:
        fail(f"GitHub API request failed ({response.status_code}) for {response.url}: {response.text}")
    return response


def iterate_github_pages(path: str, params: Dict) -> Iterable[List[Dict]]:
    page = 1
    while True:
        paged_params = dict(params)
        paged_params["page"] = page
        response = github_get(f"{GITHUB_API_URL}{path}", params=paged_params)
        items = response.json()
        if not items:
            break
        yield items
        if len(items) < int(params.get("per_page", 100)):
            break
        page += 1


def ensure_elasticsearch_ready() -> None:
    print("Waiting for Elasticsearch...")
    for _ in range(30):
        try:
            response = requests.get(f"{ELASTIC_URL}/_cluster/health", timeout=5)
            if response.status_code == 200:
                print("Elasticsearch is ready.")
                return
        except requests.RequestException:
            pass
        time.sleep(2)
    fail("Elasticsearch is not reachable at ELASTIC_URL.")


def ensure_index(index_name: str, mappings: Dict) -> None:
    response = requests.get(f"{ELASTIC_URL}/{index_name}", timeout=REQUEST_TIMEOUT_SECONDS)
    if response.status_code == 404:
        create_resp = requests.put(
            f"{ELASTIC_URL}/{index_name}",
            headers={"Content-Type": "application/json"},
            data=json.dumps(mappings),
            timeout=REQUEST_TIMEOUT_SECONDS,
        )
        if create_resp.status_code >= 400:
            fail(f"Cannot create index {index_name}: {create_resp.text}")
        print(f"Created index: {index_name}")
    elif response.status_code >= 400:
        fail(f"Cannot read index {index_name}: {response.text}")


def bulk_index(index_name: str, docs: List[Dict]) -> None:
    if not docs:
        print(f"No documents to index for {index_name}.")
        return

    lines = []
    for doc in docs:
        lines.append(json.dumps({"index": {"_index": index_name, "_id": doc["id"]}}, ensure_ascii=True))
        lines.append(json.dumps(doc, ensure_ascii=True))
    payload = "\n".join(lines) + "\n"

    response = requests.post(
        f"{ELASTIC_URL}/_bulk?refresh=true",
        headers={"Content-Type": "application/x-ndjson"},
        data=payload.encode("utf-8"),
        timeout=REQUEST_TIMEOUT_SECONDS,
    )
    if response.status_code >= 400:
        fail(f"Bulk index request failed for {index_name}: {response.text}")

    result = response.json()
    if result.get("errors"):
        fail(f"Bulk index returned errors for {index_name}: {response.text}")
    print(f"Indexed {len(docs)} documents into {index_name}.")


def issue_doc(item: Dict) -> Dict:
    labels = [label.get("name") for label in item.get("labels", []) if label.get("name")]
    assignees = [a.get("login") for a in item.get("assignees", []) if a.get("login")]
    return {
        "id": f"issue_{item['id']}",
        "repo": f"{GITHUB_OWNER}/{GITHUB_REPO}",
        "number": item.get("number"),
        "title": item.get("title"),
        "state": item.get("state"),
        "author_login": (item.get("user") or {}).get("login"),
        "labels": labels,
        "assignees": assignees,
        "comments_count": item.get("comments", 0),
        "created_at": item.get("created_at"),
        "updated_at": item.get("updated_at"),
        "closed_at": item.get("closed_at"),
        "time_to_close_hours": duration_hours(item.get("created_at"), item.get("closed_at")),
        "url": item.get("html_url"),
    }


def pull_doc(item: Dict) -> Dict:
    closed_or_merged_at = item.get("merged_at") or item.get("closed_at")
    labels = [label.get("name") for label in item.get("labels", []) if label.get("name")]
    assignees = [a.get("login") for a in item.get("assignees", []) if a.get("login")]
    reviewers = [r.get("login") for r in item.get("requested_reviewers", []) if r.get("login")]
    return {
        "id": f"pr_{item['id']}",
        "repo": f"{GITHUB_OWNER}/{GITHUB_REPO}",
        "number": item.get("number"),
        "title": item.get("title"),
        "state": item.get("state"),
        "is_draft": item.get("draft", False),
        "is_merged": item.get("merged_at") is not None,
        "author_login": (item.get("user") or {}).get("login"),
        "labels": labels,
        "assignees": assignees,
        "requested_reviewers": reviewers,
        "comments_count": item.get("comments", 0),
        "review_comments_count": item.get("review_comments", 0),
        "commits_count": item.get("commits", 0),
        "changed_files_count": item.get("changed_files", 0),
        "additions": item.get("additions", 0),
        "deletions": item.get("deletions", 0),
        "created_at": item.get("created_at"),
        "updated_at": item.get("updated_at"),
        "closed_at": item.get("closed_at"),
        "merged_at": item.get("merged_at"),
        "cycle_time_hours": duration_hours(item.get("created_at"), closed_or_merged_at),
        "url": item.get("html_url"),
    }


def commit_doc(item: Dict) -> Dict:
    commit_data = item.get("commit", {})
    author_data = commit_data.get("author", {}) or {}
    committer_data = commit_data.get("committer", {}) or {}
    return {
        "id": f"commit_{item.get('sha')}",
        "repo": f"{GITHUB_OWNER}/{GITHUB_REPO}",
        "sha": item.get("sha"),
        "author_login": (item.get("author") or {}).get("login"),
        "committer_login": (item.get("committer") or {}).get("login"),
        "author_name": author_data.get("name"),
        "author_email": author_data.get("email"),
        "authored_at": author_data.get("date"),
        "committed_at": committer_data.get("date"),
        "message": commit_data.get("message"),
        "url": item.get("html_url"),
    }


def load_issues() -> List[Dict]:
    docs: List[Dict] = []
    params = {"state": "all", "per_page": 100, "sort": "updated", "direction": "desc"}
    for page_items in iterate_github_pages(f"/repos/{GITHUB_OWNER}/{GITHUB_REPO}/issues", params):
        for item in page_items:
            if "pull_request" in item:
                continue
            docs.append(issue_doc(item))
    print(f"Fetched {len(docs)} issues.")
    return docs


def load_pull_requests() -> List[Dict]:
    docs: List[Dict] = []
    params = {"state": "all", "per_page": 100, "sort": "updated", "direction": "desc"}
    for page_items in iterate_github_pages(f"/repos/{GITHUB_OWNER}/{GITHUB_REPO}/pulls", params):
        for item in page_items:
            docs.append(pull_doc(item))
    print(f"Fetched {len(docs)} pull requests.")
    return docs


def load_commits() -> List[Dict]:
    docs: List[Dict] = []
    since = (datetime.now(timezone.utc) - timedelta(days=LOOKBACK_DAYS)).strftime("%Y-%m-%dT%H:%M:%SZ")
    params = {"per_page": 100, "since": since}
    for page_items in iterate_github_pages(f"/repos/{GITHUB_OWNER}/{GITHUB_REPO}/commits", params):
        for item in page_items:
            docs.append(commit_doc(item))
    print(f"Fetched {len(docs)} commits for last {LOOKBACK_DAYS} days.")
    return docs


def main() -> None:
    if not GITHUB_OWNER or not GITHUB_REPO:
        fail("Set GITHUB_OWNER and GITHUB_REPO environment variables.")

    if not GITHUB_TOKEN:
        print("WARNING: GITHUB_TOKEN is not set. Requests may hit strict rate limits.")

    ensure_elasticsearch_ready()

    ensure_index(
        ISSUES_INDEX,
        {
            "mappings": {
                "properties": {
                    "created_at": {"type": "date"},
                    "updated_at": {"type": "date"},
                    "closed_at": {"type": "date"},
                    "time_to_close_hours": {"type": "float"},
                }
            }
        },
    )
    ensure_index(
        PULLS_INDEX,
        {
            "mappings": {
                "properties": {
                    "created_at": {"type": "date"},
                    "updated_at": {"type": "date"},
                    "closed_at": {"type": "date"},
                    "merged_at": {"type": "date"},
                    "cycle_time_hours": {"type": "float"},
                }
            }
        },
    )
    ensure_index(
        COMMITS_INDEX,
        {
            "mappings": {
                "properties": {
                    "authored_at": {"type": "date"},
                    "committed_at": {"type": "date"},
                }
            }
        },
    )

    bulk_index(ISSUES_INDEX, load_issues())
    bulk_index(PULLS_INDEX, load_pull_requests())
    bulk_index(COMMITS_INDEX, load_commits())

    print("Done. Open Kibana at http://localhost:5601")


if __name__ == "__main__":
    main()
