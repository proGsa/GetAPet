# Генерация сваггер документации
Из папки backend выполнить в терминале команду:

 swag init -d ./cmd,./internal -g main.go --parseInternal --parseDependency --parseDepth 2 -o ./docs