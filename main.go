package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var items = []Item{
	{ID: 1, Name: "apple"},
	{ID: 2, Name: "banana"},
	{ID: 3, Name: "cherry"},
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok", Version: "1.0.0"})
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NeoFruit | Modern API Dashboard</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:opsz,wght@14..32,300;14..32,400;14..32,500;14..32,600;14..32,700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            background: linear-gradient(145deg, #f6f9fc 0%, #edf2f9 100%);
            font-family: 'Inter', sans-serif;
            color: #1a2c3e;
            line-height: 1.4;
            padding: 2rem 1.5rem;
            min-height: 100vh;
        }

        .container {
            max-width: 1280px;
            margin: 0 auto;
        }

        .header-section {
            margin-bottom: 3rem;
            display: flex;
            flex-wrap: wrap;
            justify-content: space-between;
            align-items: flex-end;
            gap: 1rem;
        }

        .title-badge h1 {
            font-size: 2.5rem;
            font-weight: 700;
            background: linear-gradient(135deg, #1E3C2C, #2A6E3A);
            background-clip: text;
            -webkit-background-clip: text;
            color: transparent;
            letter-spacing: -0.02em;
            display: inline-flex;
            align-items: center;
            gap: 0.6rem;
        }

        .title-badge h1 i {
            background: none;
            color: #2a7f3e;
            font-size: 2.2rem;
        }

        .sub {
            color: #4a627a;
            margin-top: 0.5rem;
            font-weight: 500;
            font-size: 1rem;
            border-left: 3px solid #3b9e5c;
            padding-left: 0.8rem;
        }

        .status-card {
            background: rgba(255,255,255,0.75);
            backdrop-filter: blur(12px);
            border-radius: 42px;
            padding: 0.8rem 1.6rem;
            display: flex;
            align-items: center;
            gap: 1rem;
            box-shadow: 0 10px 20px -8px rgba(0,0,0,0.08);
            border: 1px solid rgba(59,158,92,0.2);
        }

        .health-led {
            display: flex;
            align-items: center;
            gap: 12px;
        }

        .led {
            width: 14px;
            height: 14px;
            border-radius: 50%;
            background-color: #aaa;
            box-shadow: 0 0 0 2px rgba(0,0,0,0.05);
            transition: all 0.2s ease;
        }

        .led.green {
            background-color: #2bc76f;
            box-shadow: 0 0 6px #2bc76f;
        }

        .version-pill {
            background: #eef3fc;
            padding: 0.3rem 1rem;
            border-radius: 40px;
            font-size: 0.85rem;
            font-weight: 600;
            font-family: monospace;
            color: #236b3a;
        }

        .dashboard-grid {
            display: grid;
            grid-template-columns: 1fr 1.2fr;
            gap: 2rem;
            margin-bottom: 3rem;
        }

        .card {
            background: rgba(255,255,255,0.96);
            border-radius: 2rem;
            box-shadow: 0 20px 35px -12px rgba(0,0,0,0.08), 0 0 0 1px rgba(0,0,0,0.01);
            overflow: hidden;
            transition: transform 0.2s ease, box-shadow 0.2s ease;
        }

        .card:hover {
            transform: translateY(-3px);
            box-shadow: 0 28px 36px -16px rgba(0,0,0,0.12);
        }

        .card-header {
            padding: 1.4rem 1.8rem;
            border-bottom: 1px solid #eef2f5;
            display: flex;
            align-items: center;
            gap: 12px;
            background: #ffffffd9;
        }

        .card-header i {
            font-size: 1.7rem;
            color: #2a7f3e;
        }

        .card-header h2 {
            font-size: 1.45rem;
            font-weight: 600;
            letter-spacing: -0.2px;
        }

        .card-body {
            padding: 1.6rem 1.8rem;
        }

        .info-row {
            display: flex;
            justify-content: space-between;
            align-items: baseline;
            padding: 0.9rem 0;
            border-bottom: 1px dashed #e4e9ef;
        }

        .info-row:last-child {
            border-bottom: none;
        }

        .info-label {
            font-weight: 500;
            color: #5b6f86;
        }

        .info-value {
            font-weight: 600;
            font-family: monospace;
            font-size: 1.05rem;
            background: #f3f6fa;
            padding: 0.2rem 0.8rem;
            border-radius: 28px;
        }

        .items-list {
            display: flex;
            flex-direction: column;
            gap: 0.9rem;
        }

        .item-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            background: #fafcff;
            padding: 0.8rem 1.2rem;
            border-radius: 1.5rem;
            transition: all 0.2s;
            border: 1px solid #eef2f8;
        }

        .item-row:hover {
            background: white;
            border-color: #cde0d4;
            box-shadow: 0 6px 12px -8px rgba(0,0,0,0.05);
        }

        .item-id {
            font-weight: 700;
            background: #eaf4e8;
            width: 38px;
            height: 38px;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            border-radius: 40px;
            color: #1f6137;
            font-size: 1rem;
        }

        .item-name {
            font-size: 1.25rem;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 12px;
        }

        .fruit-badge {
            background: #eef3ea;
            padding: 0.25rem 0.9rem;
            border-radius: 60px;
            font-size: 0.75rem;
            font-weight: 600;
            text-transform: uppercase;
            color: #2c6e3f;
        }

        .action-bar {
            display: flex;
            justify-content: flex-end;
            margin-bottom: 1rem;
        }

        .btn-refresh {
            background: #1e2f2a;
            border: none;
            color: white;
            font-weight: 500;
            padding: 0.7rem 1.6rem;
            border-radius: 40px;
            font-family: inherit;
            display: inline-flex;
            align-items: center;
            gap: 10px;
            cursor: pointer;
            transition: all 0.2s;
        }

        .btn-refresh:hover {
            background: #14322a;
            transform: scale(0.97);
        }

        .loading-message {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 1.2rem;
            background: #eef2fc;
            border-radius: 1.5rem;
            color: #3466a1;
        }

        .error-message {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 1.2rem;
            background: #fff4e8;
            border-radius: 1.5rem;
            color: #b45f2b;
        }

        .footer-note {
            margin-top: 3.5rem;
            text-align: center;
            font-size: 0.85rem;
            color: #6c828e;
            border-top: 1px solid #dce5ec;
            padding-top: 2rem;
            display: flex;
            justify-content: center;
            gap: 2rem;
            flex-wrap: wrap;
        }

        .endpoint-badge {
            font-family: monospace;
            background: #eef0f3;
            padding: 0.2rem 1rem;
            border-radius: 30px;
            font-size: 0.8rem;
        }

        @media (max-width: 780px) {
            .dashboard-grid {
                grid-template-columns: 1fr;
            }
            .title-badge h1 {
                font-size: 1.9rem;
            }
        }
    </style>
</head>
<body>
<div class="container">
    <div class="header-section">
        <div class="title-badge">
            <h1>
                <i class="fas fa-seedling"></i> 
                NeoFruit <span style="font-weight: 400; font-size: 1.6rem;">|</span> API Studio
            </h1>
            <div class="sub">live inventory · modern web console</div>
        </div>
        <div class="status-card" id="healthWidget">
            <div class="health-led">
                <div class="led" id="ledIndicator"></div>
                <span class="glow-text" id="healthText">Checking...</span>
            </div>
            <div class="version-pill" id="versionText">v—</div>
        </div>
    </div>

    <div class="action-bar">
        <button class="btn-refresh" id="refreshBtn">
            <i class="fas fa-sync-alt"></i> Sync Inventory
        </button>
    </div>

    <div class="dashboard-grid">
        <div class="card">
            <div class="card-header">
                <i class="fas fa-heartbeat"></i>
                <h2>System Health</h2>
            </div>
            <div class="card-body">
                <div class="info-row">
                    <span class="info-label"><i class="fas fa-check-circle"></i> API Status</span>
                    <span class="info-value" id="apiStatusValue">—</span>
                </div>
                <div class="info-row">
                    <span class="info-label"><i class="fas fa-code-branch"></i> Version</span>
                    <span class="info-value" id="apiVersionValue">—</span>
                </div>
                <div class="info-row">
                    <span class="info-label"><i class="fas fa-clock"></i> Last check</span>
                    <span class="info-value" id="lastCheckTime">—</span>
                </div>
                <div class="info-row">
                    <span class="info-label"><i class="fas fa-database"></i> Items count</span>
                    <span class="info-value" id="itemsCount">0</span>
                </div>
            </div>
        </div>

        <div class="card">
            <div class="card-header">
                <i class="fas fa-fruit-apple"></i>
                <h2>Fresh Produce · Inventory</h2>
            </div>
            <div class="card-body" id="itemsContainer">
                <div class="loading-message">
                    <i class="fas fa-spinner fa-pulse"></i> loading fruits...
                </div>
            </div>
        </div>
    </div>

    <div class="footer-note">
        <div class="endpoint-badge"><i class="fas fa-link"></i> GET /health</div>
        <div class="endpoint-badge"><i class="fas fa-link"></i> GET /items</div>
        <div><i class="fas fa-lightbulb"></i> real-time data from Go backend</div>
    </div>
</div>

<script>
    const HEALTH_URL = '/health';
    const ITEMS_URL = '/items';

    async function fetchHealth() {
        try {
            const response = await fetch(HEALTH_URL);
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const data = await response.json();
            return { success: true, data };
        } catch (err) {
            return { success: false, error: err.message };
        }
    }

    async function fetchItems() {
        try {
            const response = await fetch(ITEMS_URL);
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const data = await response.json();
            return { success: true, data };
        } catch (err) {
            return { success: false, error: err.message };
        }
    }

    function updateHealthUI(healthData, errorFlag = false) {
        const ledEl = document.getElementById('ledIndicator');
        const healthTextSpan = document.getElementById('healthText');
        const apiStatusValue = document.getElementById('apiStatusValue');
        const apiVersionValue = document.getElementById('apiVersionValue');
        const versionPill = document.getElementById('versionText');

        if (errorFlag || !healthData) {
            ledEl.className = "led";
            ledEl.style.background = "#c0392b";
            healthTextSpan.innerHTML = '<i class="fas fa-exclamation-triangle"></i> Unreachable';
            if (apiStatusValue) apiStatusValue.innerHTML = '<span style="color:#c45a32;">⚠️ Connection failed</span>';
            if (apiVersionValue) apiVersionValue.innerText = "—";
            if (versionPill) versionPill.innerText = "offline";
            return;
        }

        const status = healthData.status || "unknown";
        const version = healthData.version || "?.?.?";
        
        if (status === "ok") {
            ledEl.className = "led green";
            healthTextSpan.innerHTML = '<i class="fas fa-shield-alt"></i> healthy · online';
            if (apiStatusValue) apiStatusValue.innerHTML = `<span style="color:#2a7f3e; font-weight:600;">${status.toUpperCase()} <i class="fas fa-check-circle"></i></span>`;
        } else {
            ledEl.className = "led";
            ledEl.style.background = "#e67e22";
            if (apiStatusValue) apiStatusValue.innerHTML = `<span style="color:#e67e22;">${status}</span>`;
        }
        if (apiVersionValue) apiVersionValue.innerText = version;
        if (versionPill) versionPill.innerText = `v${version}`;
    }

    function renderItems(itemsArray) {
        const container = document.getElementById('itemsContainer');
        if (!itemsArray || itemsArray.length === 0) {
            container.innerHTML = '<div class="error-message"><i class="fas fa-info-circle"></i> No items available.</div>';
            return;
        }

        const itemsListHtml = `
            <div class="items-list">
                ${itemsArray.map(item => {
                    const icon = item.name.toLowerCase() === 'banana' ? 'fas fa-leaf' : 
                                 item.name.toLowerCase() === 'cherry' ? 'fas fa-cherries' : 'fas fa-apple-alt';
                    return `
                        <div class="item-row">
                            <div class="item-name">
                                <span class="item-id">#${item.id}</span>
                                <i class="${icon}"></i>
                                <span style="text-transform: capitalize;">${item.name}</span>
                            </div>
                            <div class="fruit-badge">
                                <i class="fas fa-tag"></i> fresh
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
        `;
        container.innerHTML = itemsListHtml;
    }

    function setItemsLoading(isLoading) {
        const container = document.getElementById('itemsContainer');
        if (isLoading) {
            container.innerHTML = '<div class="loading-message"><i class="fas fa-spinner fa-pulse"></i> fetching fresh inventory...</div>';
        }
    }

    async function syncAllData() {
        setItemsLoading(true);
        
        const [healthResult, itemsResult] = await Promise.all([fetchHealth(), fetchItems()]);
        
        if (healthResult.success && healthResult.data) {
            updateHealthUI(healthResult.data, false);
        } else {
            updateHealthUI(null, true);
        }
        
        if (itemsResult.success && Array.isArray(itemsResult.data)) {
            renderItems(itemsResult.data);
            document.getElementById('itemsCount').innerText = itemsResult.data.length;
            const now = new Date();
            document.getElementById('lastCheckTime').innerText = now.toLocaleTimeString();
        } else {
            document.getElementById('itemsContainer').innerHTML = '<div class="error-message"><i class="fas fa-times-circle"></i> Failed to load items</div>';
            document.getElementById('itemsCount').innerText = '0';
        }
    }

    window.addEventListener('DOMContentLoaded', async () => {
        await syncAllData();
        setInterval(() => syncAllData(), 30000);
        document.getElementById('refreshBtn').addEventListener('click', () => syncAllData());
    });
</script>
</body>
</html>`
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/items", itemsHandler)

	log.Println(" Server starting on http://localhost:8080")
	log.Println(" Open your browser to see the modern dashboard")
	log.Fatal(http.ListenAndServe(":8080", nil))
}