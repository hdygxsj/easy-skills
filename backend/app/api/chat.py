"""Chat API (placeholder)."""
from fastapi import APIRouter

from ..schemas import ChatMessage, ChatResponse

router = APIRouter()


@router.post("/", response_model=ChatResponse)
def chat(message: ChatMessage):
    """Chat with Agent (placeholder)."""
    return ChatResponse(
        response="Agent 功能即将实现。目前请使用 API 直接操作。",
        tool_calls=None
    )
