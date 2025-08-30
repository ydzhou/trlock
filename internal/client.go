package trlock

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hekmon/transmissionrpc/v3"
)

func setupClient(hostAddr string, hostPort string) (*transmissionrpc.Client, error) {
	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%s/transmission/rpc", hostAddr, hostPort))
	if err != nil {
		return nil, err
	}
	client, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getAllTorrents(client *transmissionrpc.Client, ctx context.Context) ([]transmissionrpc.Torrent, error) {
	torrents, err := client.TorrentGetAll(ctx)
	if err != nil {
		return nil, err
	}
	return torrents, nil
}
