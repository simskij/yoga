mod cli;
mod clients;
mod config;

use clap::Parser;
use clients::github;
use tokio;

#[tokio::main]
async fn main() {
    let args = cli::Cli::parse();
    let cfg = config::load();
    
    match args.command {
        cli::Commands::Github(action) => {
            github::process(&cfg, action).await;
        }
    }
}


