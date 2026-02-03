package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

const (
	TotalKeys      = 1000000
	ConcurrentUser = 100
)

func main() {
	ctx := context.Background()

	// 1. SETUP SQLITE (Disk - WAL Mode)
	os.Remove("benchmark.db") // Clean start
	db, err := sql.Open("sqlite3", "benchmark.db?_journal_mode=WAL&_synchronous=NORMAL")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, val TEXT)"); err != nil {
		log.Fatal(err)
	}
	// Disable FS sync for seeding speed (we are benchmarking reads)
	db.Exec("PRAGMA synchronous = OFF")

	// 2. SETUP REDIS (RAM)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis not connected: %v", err)
	}
	rdb.FlushAll(ctx)

	// 3. SEED DATA
	fmt.Printf("🌱 Seeding %d keys to DB and Redis...\n", TotalKeys)
	tx, _ := db.Begin()
	pipe := rdb.Pipeline()
	for i := 0; i < TotalKeys; i++ {
		tx.Exec("INSERT INTO items (id, val) VALUES (?, ?)", i, "benchmark_value")
		pipe.Set(ctx, fmt.Sprintf("item:%d", i), "benchmark_value", 0)
		if i%1000 == 0 {
			pipe.Exec(ctx)
			pipe = rdb.Pipeline()
		}
	}
	tx.Commit()
	pipe.Exec(ctx)
	
	// 4. SETUP POSTGRES (Network SQL)
	pgConnStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	pg, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatalf("Postgres not connected: %v", err)
	}
	defer pg.Close()

	if err := pg.Ping(); err != nil {
		log.Fatalf("Postgres ping failed: %v", err)
	}
	
	if _, err := pg.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, val TEXT)"); err != nil {
		// Ignore error if table exists or connection weird, but try to continue
		log.Printf("PG Create table: %v", err)
	}
	pg.Exec("TRUNCATE TABLE items") // Clean start

	fmt.Printf("🐘 Seeding PG...\n")
	pgTx, _ := pg.Begin()
	for i := 0; i < TotalKeys; i++ {
		pgTx.Exec("INSERT INTO items (id, val) VALUES ($1, $2)", i, "benchmark_value")
	}
	pgTx.Commit()

	fmt.Println("✅ Seed Complete. Starting Benchmark...\n")

	// 5. BENCHMARK
	// A. Sequential Reads (Latency Test)
	fmt.Println("--- TEST 1: SEQUENTIAL LATENCY (1 User) ---")
	
	start := time.Now()
	for i := 0; i < TotalKeys; i++ {
		var val string
		db.QueryRow("SELECT val FROM items WHERE id = ?", i).Scan(&val)
	}
	dbDuration := time.Since(start)
	fmt.Printf("💿 SQLite (Disk/In-Process): %v (%.2f ms/op)\n", dbDuration, float64(dbDuration.Milliseconds())/float64(TotalKeys))

	start = time.Now()
	for i := 0; i < TotalKeys; i++ {
		var val string
		pg.QueryRow("SELECT val FROM items WHERE id = $1", i).Scan(&val)
	}
	pgDuration := time.Since(start)
	fmt.Printf("🐘 Postgres (Network/SQL):   %v (%.2f ms/op)\n", pgDuration, float64(pgDuration.Milliseconds())/float64(TotalKeys))

	start = time.Now()
	for i := 0; i < TotalKeys; i++ {
		rdb.Get(ctx, fmt.Sprintf("item:%d", i)).Result()
	}
	redisDuration := time.Since(start)
	fmt.Printf("⚡ Redis   (Network/RAM):   %v (%.2f ms/op)\n", redisDuration, float64(redisDuration.Milliseconds())/float64(TotalKeys))
	
	fmt.Printf("\n🚀 Latency Analysis:\n")
	fmt.Printf("- SQLite vs PG: SQLite is %.1fx faster (No Network)\n", float64(pgDuration)/float64(dbDuration))
	fmt.Printf("- PG vs Redis:  Redis is %.1fx faster (RAM vs Disk)\n\n", float64(pgDuration)/float64(redisDuration))

	// B. Concurrent Reads (Throughput Test)
	fmt.Println("--- TEST 2: CONCURRENT THROUGHPUT (100 Users) ---")
	
	measureRPS("💿 SQLite  ", func(id int) {
		var val string
		db.QueryRow("SELECT val FROM items WHERE id = ?", id).Scan(&val)
	})

	measureRPS("🐘 Postgres", func(id int) {
		var val string
		pg.QueryRow("SELECT val FROM items WHERE id = $1", id).Scan(&val)
	})

	measureRPS("⚡ Redis   ", func(id int) {
		rdb.Get(ctx, fmt.Sprintf("item:%d", id)).Result()
	})
}

func measureRPS(name string, op func(id int)) {
	var wg sync.WaitGroup
	start := time.Now()
	
	// Launch 100 concurrent workers
	workChan := make(chan int, TotalKeys)
	for i := 0; i < ConcurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range workChan {
				op(id)
			}
		}()
	}

	// Feed work
	for i := 0; i < TotalKeys; i++ {
		workChan <- i
	}
	close(workChan)
	wg.Wait()
	
	duration := time.Since(start)
	rps := float64(TotalKeys) / duration.Seconds()
	fmt.Printf("%s: %.2f RPS (Total: %v)\n", name, rps, duration)
}
