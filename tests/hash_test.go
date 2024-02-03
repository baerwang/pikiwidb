/*
 * Copyright (c) 2023-present, Qihoo, Inc.  All rights reserved.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

package pikiwidb_test

import (
	"context"
	"log"
	"strconv"
    "time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"

	"github.com/OpenAtomFoundation/pikiwidb/tests/util"
)

var _ = Describe("Hash", Ordered, func() {
	var (
		ctx    = context.TODO()
		s      *util.Server
		client *redis.Client
	)

	// BeforeAll closures will run exactly once before any of the specs
	// within the Ordered container.
	BeforeAll(func() {
		config := util.GetConfPath(false, 0)

		s = util.StartServer(config, map[string]string{"port": strconv.Itoa(7777)}, true)
		Expect(s).NotTo(Equal(nil))
	})

	// AfterAll closures will run exactly once after the last spec has
	// finished running.
	AfterAll(func() {
		err := s.Close()
		if err != nil {
			log.Println("Close Server fail.", err.Error())
			return
		}
	})

	// When running each spec Ginkgo will first run the BeforeEach
	// closure and then the subject closure.Doing so ensures that
	// each spec has a pristine, correctly initialized, copy of the
	// shared variable.
	BeforeEach(func() {
		client = s.NewClient()
		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
        time.Sleep(1 * time.Second)
	})

	// nodes that run after the spec's subject(It).
	AfterEach(func() {
		err := client.Close()
		if err != nil {
			log.Println("Close client conn fail.", err.Error())
			return
		}
	})

	//TODO(dingxiaoshuai) Add more test cases.
	It("HSet & HGet", func() {
		hSet := client.HSet(ctx, "hash", "key", "hello")
		Expect(hSet.Err()).NotTo(HaveOccurred())

		hGet := client.HGet(ctx, "hash", "key")
		Expect(hGet.Err()).NotTo(HaveOccurred())
		Expect(hGet.Val()).To(Equal("hello"))

		hGet = client.HGet(ctx, "hash", "key1")
		Expect(hGet.Err()).To(Equal(redis.Nil))
		Expect(hGet.Val()).To(Equal(""))
	})

	It("HGet & HSet 2", func() {
		testKey := "hget-hset2"
		_, err := client.Del(ctx, testKey).Result()
		Expect(err).NotTo(HaveOccurred())

		ok, err := client.HSet(ctx, testKey, map[string]interface{}{
			"key1": "hello1",
		}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(Equal(int64(1)))

		ok, err = client.HSet(ctx, testKey, map[string]interface{}{
			"key2": "hello2",
		}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(Equal(int64(1)))

		v, err := client.HGet(ctx, testKey, "key1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal("hello1"))

		v, err = client.HGet(ctx, testKey, "key2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal("hello2"))

		keys, err := client.HKeys(ctx, testKey).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(keys).To(ConsistOf([]string{"key1", "key2"}))
	})

	It("HDel", func() {
		testKey := "hdel"
		hSet := client.HSet(ctx, testKey, "key", "hello")
		Expect(hSet.Err()).NotTo(HaveOccurred())
		hSet = client.HSet(ctx, testKey, "key", "hello")
		Expect(hSet.Err()).NotTo(HaveOccurred())
		Expect(hSet.Val()).To(Equal(int64(0)))

		hDel := client.HDel(ctx, testKey, "key")
		Expect(hDel.Err()).NotTo(HaveOccurred())
		Expect(hDel.Val()).To(Equal(int64(1)))

		hDel = client.HDel(ctx, testKey, "key")
		Expect(hDel.Err()).NotTo(HaveOccurred())
		Expect(hDel.Val()).To(Equal(int64(0)))

		hSet = client.HSet(ctx, testKey, "key", "hello")
		Expect(hSet.Err()).NotTo(HaveOccurred())
		Expect(hSet.Val()).To(Equal(int64(1)))

		hDel = client.HDel(ctx, testKey, "key")
		Expect(hDel.Err()).NotTo(HaveOccurred())
		Expect(hDel.Val()).To(Equal(int64(1)))
	})

	It("HGetAll", func() {
		testKey := "hgetall"
		err := client.HSet(ctx, testKey, "key1", "hello1").Err()
		Expect(err).NotTo(HaveOccurred())
		err = client.HSet(ctx, testKey, "key2", "hello2").Err()
		Expect(err).NotTo(HaveOccurred())

		m, err := client.HGetAll(ctx, testKey).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(m).To(Equal(map[string]string{"key1": "hello1", "key2": "hello2"}))
	})

	It("HMGet", func() {
		testKey := "hmget"
		err := client.HSet(ctx, testKey, "key1", "hello1").Err()
		Expect(err).NotTo(HaveOccurred())

		vals, err := client.HMGet(ctx, testKey, "key1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(vals).To(Equal([]interface{}{"hello1"}))
	})

	It("HKeys", func() {
		testKey := "hkeys"
		hkeys := client.HKeys(ctx, testKey)
		Expect(hkeys.Err()).NotTo(HaveOccurred())
		Expect(hkeys.Val()).To(Equal([]string{}))

		hset := client.HSet(ctx, testKey, "key1", "hello1")
		Expect(hset.Err()).NotTo(HaveOccurred())
		hset = client.HSet(ctx, testKey, "key2", "hello2")
		Expect(hset.Err()).NotTo(HaveOccurred())

		hkeys = client.HKeys(ctx, testKey)
		Expect(hkeys.Err()).NotTo(HaveOccurred())
		Expect(hkeys.Val()).To(Equal([]string{"key1", "key2"}))
	})

	It("HLen", func() {
		testKey := "hlen"
		hSet := client.HSet(ctx, testKey, "key1", "hello1")
		Expect(hSet.Err()).NotTo(HaveOccurred())
		hSet = client.HSet(ctx, testKey, "key2", "hello2")
		Expect(hSet.Err()).NotTo(HaveOccurred())

		hLen := client.HLen(ctx, testKey)
		Expect(hLen.Err()).NotTo(HaveOccurred())
		Expect(hLen.Val()).To(Equal(int64(2)))
	})

	It("HStrLen", func() {
		testKey := "hstrlen"
		hSet := client.HSet(ctx, testKey, "key1", "hello1")
		Expect(hSet.Err()).NotTo(HaveOccurred())

		hGet := client.HGet(ctx, testKey, "key1")
		Expect(hGet.Err()).NotTo(HaveOccurred())
		length := client.Do(ctx, "hstrlen", testKey, "key1")

		Expect(length.Val()).To(Equal(int64(len("hello1"))))
	})

})