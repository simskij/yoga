use std::env;
use serde::{Deserialize, Serialize};
use std::path::PathBuf;

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Config {
    pub github: GitHub
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct GitHub {
    pub token: Option<String>,
    pub repos: Vec<String>,
}

impl GitHub {
    pub fn token(&self) -> String {
        let key = "GITHUB_TOKEN";
        match env::var_os(key) {
            Some(val) => val.into_string().unwrap()  ,
            None => {
                self.token.clone().unwrap_or_else(|| String::from(""))
            }
        }
    }

}

pub fn load() -> Config {
    let home = dirs::home_dir().unwrap_or(PathBuf::from(""));
    let file = format!("{}/.config/yoga/config.yaml", home.display());
    
    let content: String = std::fs::read_to_string(file).expect("Unable to open file");
    serde_yaml::from_str(&content).unwrap()
}
