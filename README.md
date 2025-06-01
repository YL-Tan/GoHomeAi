# GoHomeAi

**A Golang-powered, AI-driven smart-home automation platform for real-time monitoring, automated control, and modular hardware integration.**

## About

GoHomeAi is designed as a modular, full-stack platform that combines:

- **Golang** for high-performance backend microservices  
- **FastAPI** (Python) for lightweight AI endpoints  
- **Nuxt 3 + Tailwind CSS** for the user-facing dashboard  
- **PyTorch + scikit-learn** for machine-learning models  
- **PostgreSQL** for persistent storage of device data and telemetry  
- **Docker Compose** for containerized local development  
- **MLflow** (or DVC) to track model experiments and deployments

The goal is to provide a robust software foundation for:

1. **Real-time Monitoring**  
   - Collect and store telemetry (energy usage, sensor readings)  
   - Expose REST/WebSocket endpoints for live updates  

2. **Automated Control**  
   - Define rule-based or AI-driven automation (e.g., thermostat / lighting adjustments)  
   - Extendable to plug-in hardware controllers (Zigbee, Z-Wave, Matter)

3. **Modular Hardware Integration**  
   - Easily swap in/out real devices (Raspberry Pi, ESP32, smart plugs)  
   - Simulate data if hardware is not yet available

## Prerequisites

Before you begin, ensure you have the following installed on your local machine:

1. TBD

> **Note:** The instructions above assume you are running on a Unix-like OS (Linux/macOS). Adjust accordingly for Windows.

---
## Getting Started

### Clone the Repository

```bash
git clone https://github.com/YL-Tan/GoHomeAi.git
cd gohomeai
```

