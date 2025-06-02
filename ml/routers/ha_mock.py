from fastapi import APIRouter
from pydantic import BaseModel

router = APIRouter(prefix="/ha")


class StatusResponse(BaseModel):
    status: str
    uptime: float


@router.get("/status", response_model=StatusResponse)
async def get_status():
    return {"status": "online", "uptime": 12345.6}
