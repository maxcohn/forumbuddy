use actix_web::web::Data;
use sqlx::postgres::PgPoolOptions;
use sqlx::postgres::PgConnectOptions;
use sqlx::{Pool, Postgres};
use sqlx::types::chrono::{DateTime, Utc};
use sqlx;
use anyhow::{Result, anyhow, bail};
use actix_web::{get, web, HttpServer, Responder, App};
use serde::{Serialize, Deserialize};
use std::fmt;

#[derive(Debug, Serialize, Deserialize)]
struct User {
    pub uid: i32,
    pub username: String,
    //#[serde(with="chrono::serde::ts_milliseconds")]
    pub created_at: DateTime<Utc>,
}

#[derive(Debug)]
struct WebError(anyhow::Error);

impl fmt::Display for WebError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl actix_web::error::ResponseError for WebError {

}

impl From<anyhow::Error> for WebError {
    
}


impl User {
    pub async fn get_by_id(pool: &Pool<Postgres>, id: i64) -> Result<Self> {
        let mut potential_users = sqlx::query_as!(User, "SELECT * FROM users WHERE uid = $1 LIMIT 1", id as i32)
            .fetch_all(pool)
            .await?;

        match potential_users.pop() {
            Some(user) => Ok(user),
            None => bail!("Failed to find user with id: {}", id)
        }
    }
}


#[get("/user/{id}")]
async fn get_user(app_state: web::Data<AppState>, id: web::Path<i64>) -> impl Responder {
    //format!("path id: {:?}", User::get_by_id(&app_state.db, *id).await.unwrap())

    let user = match User::get_by_id(&app_state.db, *id).await {
        Ok(u) => u,
        Err(e) => return Err(())
    };

    Ok(web::Json(user))
}

#[derive(Debug, Clone)]
struct AppState {
    pub db: Pool<Postgres>
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let pool = PgPoolOptions::new()
        .max_connections(20)
        .connect_with(PgConnectOptions::new()
            .port(5432)
            .host("127.0.0.1")
            .password("password")
            .username("postgres")).await.unwrap();

    let app_state = AppState {
        db: pool
    };
    
    HttpServer::new(move || {
        App::new()
            .app_data(Data::new(app_state.clone()))
            .service(get_user)
    })
        .bind(("127.0.0.1", 8080))?
        .run()
        .await

}
