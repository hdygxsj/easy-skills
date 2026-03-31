"""API routes."""
from fastapi import APIRouter
from .packages import router as packages_router
from .groups import router as groups_router
from .installations import router as installations_router
from .config import router as config_router
from .chat import router as chat_router

api_router = APIRouter()

api_router.include_router(packages_router, prefix="/packages", tags=["packages"])
api_router.include_router(groups_router, prefix="/groups", tags=["groups"])
api_router.include_router(installations_router, prefix="/installations", tags=["installations"])
api_router.include_router(config_router, prefix="/config", tags=["config"])
api_router.include_router(chat_router, prefix="/chat", tags=["chat"])
