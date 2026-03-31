"""Group management API."""
from typing import List
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from ..database import get_db
from ..schemas import GroupCreate, GroupResponse, GroupDetail, PackageResponse
from ..services import GroupService

router = APIRouter()


@router.get("/", response_model=List[GroupResponse])
def list_groups(db: Session = Depends(get_db)):
    """List all groups."""
    service = GroupService(db)
    return service.list_groups()


@router.get("/{group_id}", response_model=GroupDetail)
def get_group(group_id: int, db: Session = Depends(get_db)):
    """Get group details with packages."""
    service = GroupService(db)
    group = service.get_group(group_id)
    if not group:
        raise HTTPException(status_code=404, detail="Group not found")
    return group


@router.post("/", response_model=GroupResponse)
def create_group(data: GroupCreate, db: Session = Depends(get_db)):
    """Create a new group."""
    service = GroupService(db)
    try:
        return service.create_group(data.name, data.description)
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.delete("/{group_id}")
def delete_group(group_id: int, db: Session = Depends(get_db)):
    """Delete a group."""
    service = GroupService(db)
    if not service.delete_group(group_id):
        raise HTTPException(status_code=404, detail="Group not found")
    return {"status": "deleted"}


@router.post("/{group_id}/packages/{package_id}")
def add_package_to_group(group_id: int, package_id: int, db: Session = Depends(get_db)):
    """Add package to group."""
    service = GroupService(db)
    if not service.add_package_to_group(group_id, package_id):
        raise HTTPException(status_code=404, detail="Group or package not found")
    return {"status": "added"}


@router.delete("/{group_id}/packages/{package_id}")
def remove_package_from_group(group_id: int, package_id: int, db: Session = Depends(get_db)):
    """Remove package from group."""
    service = GroupService(db)
    if not service.remove_package_from_group(group_id, package_id):
        raise HTTPException(status_code=404, detail="Package not in group")
    return {"status": "removed"}
