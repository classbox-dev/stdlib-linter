package main

import (
	"context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type args struct {
	Root string `positional-arg-name:"path"`
}

type Options struct {
	ConfigUrl  string `short:"u" env:"CONFIG_URL" description:"URL with remote JSON config" default:"https://hsecode.com/stdlib/docs/linter.yaml"`
	ConfigPath string `short:"c" env:"CONFIG_PATH" description:"path to local config that will be used instead of remote one"`
	Args       args   `positional-args:"args"`
	Verbose    bool   `short:"v" description:"print more errors"`
}

func (opts *Options) GetConfig() *Config {
	var config Config
	if opts.ConfigPath != "" {
		content, err := ioutil.ReadFile(opts.ConfigPath)
		if err != nil {
			log.Fatalf("could not open %v: %v", opts.ConfigPath, err)
		}
		if err := yaml.Unmarshal(content, &config); err != nil {
			log.Fatalf("could not parse %v: %v", opts.ConfigPath, err)
		}
		return &config
	}
	err := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", opts.ConfigUrl, nil)
		if err != nil {
			return err
		}
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		//noinspection GoUnhandledErrorResult
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(content, &config); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		log.Printf("[WARN] could not read remote config")
		if opts.Verbose {
			log.Print(err)
		}
		log.Print("[WARN] using default config (may be incomplete)")
		return &defaultConfig
	}
	return &config
}
