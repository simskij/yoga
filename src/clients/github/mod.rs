use tabled::Table;
use tabled::Tabled;
use crate::cli::GithubArgs;
use crate::{cli};
use crate::config::Config;

pub async fn list_issue_counts(_cfg: &Config) -> Result<(), octocrab::Error> {
    let mut builder = octocrab::Octocrab::builder();
    let token = _cfg.github.token();

    if token.is_empty() {
        println!("Warning: No GITHUB_TOKEN set. Will run anonymously, which might incur rate limits.");
    } else {
        builder = builder.personal_token(token);
    }

    let crab = builder.build()?;
    let mut handles = vec![];

    for name in _cfg.github.repos.clone() {
        let handle = async {
            let issues = get_repo_counts(&crab, &name, "issue").await;
            let prs = get_repo_counts(&crab, &name, "pr").await;

            RepoStats { name, issues, prs }
        };
        handles.push(handle);
    }

    let repos = futures::future::join_all(handles).await;

    let output = Table::new(repos)
        .with(tabled::settings::Style::sharp())
        .to_string();

    println!("{}", output);
    Ok(())
}

async fn get_repo_counts(crab: &octocrab::Octocrab, repo: &String, kind: &str) -> u64 {
    let res: Result<serde_json::Value, octocrab::Error> = crab.graphql(
            &serde_json::json!({
                "query": format!(
                    "{{ search(type: ISSUE, query: \"repo:{} is:open is:{}\") {{ issueCount }} }}",
                    repo,
                    kind
                )
            })
        )
        .await;

    match res {
        Ok(res) => res["data"]["search"]["issueCount"]
                .as_u64()
                .unwrap()
                .clone(),
        Err(_error) => 0
    }
}

#[derive(Tabled, Debug)]
struct RepoStats {
    #[tabled(rename = "Repository")]
    name: String,
    #[tabled(rename = "Issues")]
    issues: u64,
    #[tabled(rename = "PRs")]
    prs: u64,
}

pub async fn process(config: &Config, args: GithubArgs) {
    let gh_cmd = args.command.unwrap();
    match gh_cmd {
        cli::GitHubCommands::Issues => {
            list_issue_counts(config)
                .await.expect("Could not load issue counts.");
        }
    }
}