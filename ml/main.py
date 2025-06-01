from typing import List, Optional

from fastapi import FastAPI
from pydantic import BaseModel, Field
from routers.health import router as health_router

api = FastAPI(title="GoHomeAI-Inference Service")

api.include_router(health_router)
