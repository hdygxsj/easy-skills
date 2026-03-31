"""Tests for package management."""
import pytest


def test_list_packages_empty(client):
    """Test listing packages when empty."""
    response = client.get("/api/packages/")
    assert response.status_code == 200
    assert response.json() == []


def test_create_group(client):
    """Test creating a group."""
    response = client.post("/api/groups/", json={
        "name": "test-group",
        "description": "Test group"
    })
    assert response.status_code == 200
    data = response.json()
    assert data["name"] == "test-group"
    assert data["description"] == "Test group"


def test_list_groups(client):
    """Test listing groups."""
    # Create a group first
    client.post("/api/groups/", json={"name": "group1"})
    
    response = client.get("/api/groups/")
    assert response.status_code == 200
    assert len(response.json()) == 1


def test_delete_group(client):
    """Test deleting a group."""
    # Create a group
    create_response = client.post("/api/groups/", json={"name": "to-delete"})
    group_id = create_response.json()["id"]
    
    # Delete it
    delete_response = client.delete(f"/api/groups/{group_id}")
    assert delete_response.status_code == 200
    
    # Verify it's gone
    list_response = client.get("/api/groups/")
    assert len(list_response.json()) == 0


def test_health_check(client):
    """Test health check endpoint."""
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "healthy"
