use sqlx::postgres::PgPoolOptions;
use sqlx::postgres::PgConnectOptions;
use sqlx::types::time::OffsetDateTime;
use sqlx;


#[derive(Debug)]
struct User {
    uid: i32,
    username: String,
    password_hash: String,
    created_at: OffsetDateTime,
}
#[tokio::main]
async fn main() {
    let pool = PgPoolOptions::new()
        .max_connections(20)
        .connect_with(PgConnectOptions::new()
            .port(5432)
            .host("127.0.0.1")
            .password("password")
            .username("postgres")).await.unwrap();


    let u = sqlx::query_as!(User, "SELECT * FROM users").fetch_all(&pool).await.unwrap();

    
    println!("Hello, {:?}!", &u);
}
