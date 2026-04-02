#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use std::process::{Command, Child};
use std::thread;
use std::time::Duration;
use std::sync::{Arc, Mutex};
use tauri::Manager;

fn main() {
    // Get the app resources path
    let resources_path = std::env::current_exe()
        .unwrap()
        .parent()
        .unwrap()
        .parent()
        .unwrap()
        .join("Resources");
    
    let index_path = resources_path.join("index.html");
    let index_url = format!("file://{}", index_path.to_string_lossy());
    
    println!("Loading frontend from: {}", index_url);

    tauri::Builder::default()
        .setup(move |app| {
            let window = app.get_webview_window("main").unwrap();
            
            // Navigate to the local index.html
            window.eval(&format!("window.location.href = '{}';", index_url)).ok();
            
            // Start easy-skills serve in background for API calls
            let child = Command::new("easy-skills")
                .arg("serve")
                .spawn();
            
            let backend = Arc::new(Mutex::new(child.ok()));
            let backend_clone = backend.clone();
            
            window.on_window_event(move |event| {
                if let tauri::WindowEvent::CloseRequested { .. } = event {
                    if let Ok(mut guard) = backend_clone.lock() {
                        if let Some(ref mut c) = *guard {
                            let _ = c.kill();
                        }
                    }
                }
            });
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
