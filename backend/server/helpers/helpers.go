// helpers/helpers.go
package helpers

import (
	pbFlow "github.com/cilium/cilium/api/v1/flow"
	"github.com/cilium/hubble-ui/backend/domain/cache"
	"github.com/cilium/hubble-ui/backend/domain/flow"
	"github.com/cilium/hubble-ui/backend/domain/link"
	"github.com/cilium/hubble-ui/backend/domain/service"
	"github.com/cilium/hubble-ui/backend/internal/types"
	"github.com/cilium/hubble-ui/backend/proto/ui"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventResponseForService(
	svc *service.Service, cflags *cache.Flags,
) *ui.GetEventsResponse {
	sstate := &ui.ServiceState{
		Service: svc.ToProto(),
		Type:    StateChangeFromCacheFlags(cflags),
	}

	f := svc.FlowRef()

	return &ui.GetEventsResponse{
		Node:      f.NodeName,
		Timestamp: f.Time,
		Event: &ui.GetEventsResponse_ServiceState{
			ServiceState: sstate,
		},
	}
}

func EventResponseForLink(
	l *link.Link, cflags *cache.Flags,
) *ui.GetEventsResponse {
	f := l.IntoFlow()
	lstate := &ui.ServiceLinkState{
		ServiceLink: l.ToProto(),
		Type:        StateChangeFromCacheFlags(cflags),
	}

	return &ui.GetEventsResponse{
		Node:      f.NodeName,
		Timestamp: f.Time,
		Event: &ui.GetEventsResponse_ServiceLinkState{
			ServiceLinkState: lstate,
		},
	}
}

func EventResponseFromFlow(f *flow.Flow) *ui.GetEventsResponse {
	ref := f.Ref()
	return &ui.GetEventsResponse{
		Node:      ref.NodeName,
		Timestamp: ref.Time,
		Event: &ui.GetEventsResponse_Flow{
			Flow: ref,
		},
	}
}

func EventResponseFromRawFlows(flows []*pbFlow.Flow) *ui.GetEventsResponse {
	n := len(flows)
	if n == 0 {
		return nil
	}

	ref := flows[n-1]
	return &ui.GetEventsResponse{
		Node:      ref.NodeName,
		Timestamp: ref.Time,
		Event: &ui.GetEventsResponse_Flows{
			Flows: &ui.Flows{
				Flows: flows,
			},
		},
	}
}

func notificationConnState(
	connected,
	reconnecting,
	k8sUnavailable,
	k8sConnected bool,
) *ui.Notification {
	return &ui.Notification{
		Notification: &ui.Notification_ConnState{
			ConnState: &ui.ConnectionState{
				Connected:      connected,
				Reconnecting:   reconnecting,
				K8SUnavailable: k8sUnavailable,
				K8SConnected:   k8sConnected,
			},
		},
	}
}

func EventResponseReconnecting() *ui.GetEventsResponse {
	return &ui.GetEventsResponse{
		Node:      "backend",
		Timestamp: timestamppb.Now(),
		Event: &ui.GetEventsResponse_Notification{
			Notification: notificationConnState(false, true, false, false),
		},
	}
}

func EventResponseConnected() *ui.GetEventsResponse {
	return &ui.GetEventsResponse{
		Node:      "backend",
		Timestamp: timestamppb.Now(),
		Event: &ui.GetEventsResponse_Notification{
			Notification: notificationConnState(true, false, false, false),
		},
	}
}

func EventResponseK8sUnavailable() *ui.GetEventsResponse {
	return &ui.GetEventsResponse{
		Node:      "backend",
		Timestamp: timestamppb.Now(),
		Event: &ui.GetEventsResponse_Notification{
			Notification: notificationConnState(false, false, true, false),
		},
	}
}

func EventResponseK8sConnected() *ui.GetEventsResponse {
	return &ui.GetEventsResponse{
		Node:      "backend",
		Timestamp: timestamppb.Now(),
		Event: &ui.GetEventsResponse_Notification{
			Notification: notificationConnState(false, false, false, true),
		},
	}
}

func EventResponseNoPermission(resource string, err string) *ui.GetEventsResponse {
	notif := &ui.Notification{
		Notification: &ui.Notification_NoPermission{
			NoPermission: &ui.NoPermission{
				Resource: resource,
				Error:    err,
			},
		},
	}

	return &ui.GetEventsResponse{
		Node:      "backend",
		Timestamp: timestamppb.Now(),
		Event: &ui.GetEventsResponse_Notification{
			Notification: notif,
		},
	}
}

func StateChangeFromCacheFlags(cflags *cache.Flags) ui.StateChange {
	if cflags.Exists {
		return ui.StateChange_EXISTS
	} else if cflags.Created {
		return ui.StateChange_ADDED
	} else if cflags.Updated {
		return ui.StateChange_MODIFIED
	} else if cflags.Deleted {
		return ui.StateChange_DELETED
	}

	return ui.StateChange_UNKNOWN_STATE_CHANGE
}

func StateChangeFromEventType(evtType types.EventKind) ui.StateChange {
	if evtType == types.Added {
		return ui.StateChange_ADDED
	} else if evtType == types.Deleted {
		return ui.StateChange_DELETED
	} else if evtType == types.Modified {
		return ui.StateChange_MODIFIED
	} else if evtType == types.Exists {
		return ui.StateChange_EXISTS
	}

	return ui.StateChange_UNKNOWN_STATE_CHANGE
}
