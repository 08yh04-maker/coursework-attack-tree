# Архитектура курсовой работы

## Общая архитектура

```mermaid
graph TB
    subgraph "Пользователь"
        UI[Веб-интерфейс]
    end

    subgraph "Бэкенд Python"
        API[FastAPI]
        ANALYZER[NetworkX Analyzer]
        METRICS[Metrics Calculator]
    end

    subgraph "Go модуль"
        GO_CALC[Risk Calculator]
    end

    subgraph "Хранилище"
        JSON[(attack_trees.json)]
    end

    UI --> API
    API --> ANALYZER
    API --> GO_CALC
    ANALYZER --> JSON
    GO_CALC --> METRICS
    METRICS --> API

Компоненты системы
1. Python бэкенд (FastAPI)
API эндпоинты — создание, загрузка, анализ деревьев атак

NetworkX анализатор — построение графов, поиск путей

Метрики — вероятность успеха, стоимость атаки

2. Go модуль
Высокопроизводительный расчёт рисков

Оптимизация для больших графов

3. Веб-интерфейс
Визуализация дерева атак (интерактивный граф)

Панель управления

Отображение результатов анализа

Поток данных
flowchart LR
    User[Пользователь] -->|JSON| Create[Создание дерева]
    Create --> Validate[Валидация]
    Validate --> Graph[Построение графа]
    Graph --> GoCalc[Go: расчёт рисков]
    GoCalc --> Metrics[Метрики]
    Metrics --> Visual[Визуализация]
    Visual --> User

Формат входных данных (JSON)
{
  "name": "Атака на веб-приложение",
  "root": {
    "id": "root",
    "type": "OR",
    "children": [
      {"id": "sql_injection", "probability": 0.3, "cost": 5000},
      {"id": "xss", "probability": 0.4, "cost": 3000}
    ]
  }
}
Метрики риска
Вероятность успеха — рассчитывается рекурсивно по дереву

Стоимость реализации — суммарные затраты злоумышленника

Критичность — комбинированная метрика

Технологии
Компонент	Технология
API	FastAPI + Pydantic
Графы	NetworkX
Расчёты	Go (goroutines)
Визуализация	D3.js / Cytoscape.js
CI/CD	GitHub Actions
Контейнеризация	Docker Compose




