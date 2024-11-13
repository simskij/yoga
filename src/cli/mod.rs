use clap::Args;
use clap::Subcommand;
use clap::Parser;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
#[command(propagate_version = true)]
pub struct Cli {
    #[command(subcommand)]
    pub(crate) command: Commands,
}

#[derive(Debug, Subcommand)]
pub enum Commands {
    Github(GithubArgs),
}

#[derive(Debug, Args)]
#[command(args_conflicts_with_subcommands = true)]
pub struct GithubArgs {
    #[command(subcommand)]
    pub(crate) command: Option<GitHubCommands>,
}

#[derive(Debug, Subcommand)]
#[command(arg_required_else_help = true)]
pub enum GitHubCommands {
    Issues ,
}
