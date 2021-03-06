package nftables_exporter

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"sync"
)

var linkList = map[uint32]netlink.Link{}
var linkListLock = &sync.Mutex{}

func GetInterfaceFromNumber(number uint32) (netlink.Link, error) {
	defer linkListLock.Unlock()
	linkListLock.Lock()

	if link, ok := linkList[number]; ok {
		return link, nil
	}

	if newLinkList, err := netlink.LinkList(); err != nil {
		return nil, err
	} else {
		for _, link := range newLinkList {
			linkList[uint32(link.Attrs().Index)] = link
		}
	}

	if link, ok := linkList[number]; ok {
		return link, nil
	} else {
		return nil, fmt.Errorf("link with index %d not found", number)
	}
}
