package repository

// TODO: write Redis tests

//var redisP *RedisConnection

// func SetupTestRedis() (*rds.Client, func(), error) {
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	resource, err := pool.Run("redis", "latest", nil)
// 	if err != nil {
// 		log.Fatalf("Could not start resource: %s", err)
// 	}

// 	// if run with docker-machine the hostname needs to be set
// 	opt, err := rds.ParseURL("redis://:@localhost:6379/1")
// 	if err != nil {
// 		log.Fatalf("Could not parse endpoint: %s", pool.Client.Endpoint())
// 	}
// 	var rdb *rds.Client
// 	if err := pool.Retry(func() error {
// 		rdb := rds.NewClient(opt)

// 		ping := rdb.Ping(context.Background())
// 		return ping.Err()
// 	}); err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	cleanup := func() {
// 		rdb.Close()
// 		pool.Purge(resource)
// 	}
// 	return rdb, cleanup, nil
// }

// func TestRedisSetByID(t *testing.T) {
// 	err := redisP.RedisSetByID(context.Background(), &entityEugen)
// 	require.NoError(t, err)
// }
