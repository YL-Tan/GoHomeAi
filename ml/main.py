from typing import List, Optional

from fastapi import FastAPI
from pydantic import BaseModel, Field
from routers.ha_mock import router as ha_router
from routers.health import router as health_router

api = FastAPI(title="GoHomeAI-Inference Service")

api.include_router(health_router)
api.include_router(ha_router)
