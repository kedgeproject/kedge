// +build !ignore_autogenerated_openshift

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1

import (
	route "github.com/openshift/origin/pkg/route/apis/route"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	api "k8s.io/kubernetes/pkg/api"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	unsafe "unsafe"
)

func init() {
	SchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1_Route_To_route_Route,
		Convert_route_Route_To_v1_Route,
		Convert_v1_RouteIngress_To_route_RouteIngress,
		Convert_route_RouteIngress_To_v1_RouteIngress,
		Convert_v1_RouteIngressCondition_To_route_RouteIngressCondition,
		Convert_route_RouteIngressCondition_To_v1_RouteIngressCondition,
		Convert_v1_RouteList_To_route_RouteList,
		Convert_route_RouteList_To_v1_RouteList,
		Convert_v1_RoutePort_To_route_RoutePort,
		Convert_route_RoutePort_To_v1_RoutePort,
		Convert_v1_RouteSpec_To_route_RouteSpec,
		Convert_route_RouteSpec_To_v1_RouteSpec,
		Convert_v1_RouteStatus_To_route_RouteStatus,
		Convert_route_RouteStatus_To_v1_RouteStatus,
		Convert_v1_RouteTargetReference_To_route_RouteTargetReference,
		Convert_route_RouteTargetReference_To_v1_RouteTargetReference,
		Convert_v1_RouterShard_To_route_RouterShard,
		Convert_route_RouterShard_To_v1_RouterShard,
		Convert_v1_TLSConfig_To_route_TLSConfig,
		Convert_route_TLSConfig_To_v1_TLSConfig,
	)
}

func autoConvert_v1_Route_To_route_Route(in *Route, out *route.Route, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_RouteSpec_To_route_RouteSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1_RouteStatus_To_route_RouteStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1_Route_To_route_Route(in *Route, out *route.Route, s conversion.Scope) error {
	return autoConvert_v1_Route_To_route_Route(in, out, s)
}

func autoConvert_route_Route_To_v1_Route(in *route.Route, out *Route, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_route_RouteSpec_To_v1_RouteSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_route_RouteStatus_To_v1_RouteStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func Convert_route_Route_To_v1_Route(in *route.Route, out *Route, s conversion.Scope) error {
	return autoConvert_route_Route_To_v1_Route(in, out, s)
}

func autoConvert_v1_RouteIngress_To_route_RouteIngress(in *RouteIngress, out *route.RouteIngress, s conversion.Scope) error {
	out.Host = in.Host
	out.RouterName = in.RouterName
	out.Conditions = *(*[]route.RouteIngressCondition)(unsafe.Pointer(&in.Conditions))
	out.WildcardPolicy = route.WildcardPolicyType(in.WildcardPolicy)
	out.RouterCanonicalHostname = in.RouterCanonicalHostname
	return nil
}

func Convert_v1_RouteIngress_To_route_RouteIngress(in *RouteIngress, out *route.RouteIngress, s conversion.Scope) error {
	return autoConvert_v1_RouteIngress_To_route_RouteIngress(in, out, s)
}

func autoConvert_route_RouteIngress_To_v1_RouteIngress(in *route.RouteIngress, out *RouteIngress, s conversion.Scope) error {
	out.Host = in.Host
	out.RouterName = in.RouterName
	out.Conditions = *(*[]RouteIngressCondition)(unsafe.Pointer(&in.Conditions))
	out.WildcardPolicy = WildcardPolicyType(in.WildcardPolicy)
	out.RouterCanonicalHostname = in.RouterCanonicalHostname
	return nil
}

func Convert_route_RouteIngress_To_v1_RouteIngress(in *route.RouteIngress, out *RouteIngress, s conversion.Scope) error {
	return autoConvert_route_RouteIngress_To_v1_RouteIngress(in, out, s)
}

func autoConvert_v1_RouteIngressCondition_To_route_RouteIngressCondition(in *RouteIngressCondition, out *route.RouteIngressCondition, s conversion.Scope) error {
	out.Type = route.RouteIngressConditionType(in.Type)
	out.Status = api.ConditionStatus(in.Status)
	out.Reason = in.Reason
	out.Message = in.Message
	out.LastTransitionTime = (*meta_v1.Time)(unsafe.Pointer(in.LastTransitionTime))
	return nil
}

func Convert_v1_RouteIngressCondition_To_route_RouteIngressCondition(in *RouteIngressCondition, out *route.RouteIngressCondition, s conversion.Scope) error {
	return autoConvert_v1_RouteIngressCondition_To_route_RouteIngressCondition(in, out, s)
}

func autoConvert_route_RouteIngressCondition_To_v1_RouteIngressCondition(in *route.RouteIngressCondition, out *RouteIngressCondition, s conversion.Scope) error {
	out.Type = RouteIngressConditionType(in.Type)
	out.Status = api_v1.ConditionStatus(in.Status)
	out.Reason = in.Reason
	out.Message = in.Message
	out.LastTransitionTime = (*meta_v1.Time)(unsafe.Pointer(in.LastTransitionTime))
	return nil
}

func Convert_route_RouteIngressCondition_To_v1_RouteIngressCondition(in *route.RouteIngressCondition, out *RouteIngressCondition, s conversion.Scope) error {
	return autoConvert_route_RouteIngressCondition_To_v1_RouteIngressCondition(in, out, s)
}

func autoConvert_v1_RouteList_To_route_RouteList(in *RouteList, out *route.RouteList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]route.Route)(unsafe.Pointer(&in.Items))
	return nil
}

func Convert_v1_RouteList_To_route_RouteList(in *RouteList, out *route.RouteList, s conversion.Scope) error {
	return autoConvert_v1_RouteList_To_route_RouteList(in, out, s)
}

func autoConvert_route_RouteList_To_v1_RouteList(in *route.RouteList, out *RouteList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items == nil {
		out.Items = make([]Route, 0)
	} else {
		out.Items = *(*[]Route)(unsafe.Pointer(&in.Items))
	}
	return nil
}

func Convert_route_RouteList_To_v1_RouteList(in *route.RouteList, out *RouteList, s conversion.Scope) error {
	return autoConvert_route_RouteList_To_v1_RouteList(in, out, s)
}

func autoConvert_v1_RoutePort_To_route_RoutePort(in *RoutePort, out *route.RoutePort, s conversion.Scope) error {
	out.TargetPort = in.TargetPort
	return nil
}

func Convert_v1_RoutePort_To_route_RoutePort(in *RoutePort, out *route.RoutePort, s conversion.Scope) error {
	return autoConvert_v1_RoutePort_To_route_RoutePort(in, out, s)
}

func autoConvert_route_RoutePort_To_v1_RoutePort(in *route.RoutePort, out *RoutePort, s conversion.Scope) error {
	out.TargetPort = in.TargetPort
	return nil
}

func Convert_route_RoutePort_To_v1_RoutePort(in *route.RoutePort, out *RoutePort, s conversion.Scope) error {
	return autoConvert_route_RoutePort_To_v1_RoutePort(in, out, s)
}

func autoConvert_v1_RouteSpec_To_route_RouteSpec(in *RouteSpec, out *route.RouteSpec, s conversion.Scope) error {
	out.Host = in.Host
	out.Path = in.Path
	if err := Convert_v1_RouteTargetReference_To_route_RouteTargetReference(&in.To, &out.To, s); err != nil {
		return err
	}
	out.AlternateBackends = *(*[]route.RouteTargetReference)(unsafe.Pointer(&in.AlternateBackends))
	out.Port = (*route.RoutePort)(unsafe.Pointer(in.Port))
	out.TLS = (*route.TLSConfig)(unsafe.Pointer(in.TLS))
	out.WildcardPolicy = route.WildcardPolicyType(in.WildcardPolicy)
	return nil
}

func Convert_v1_RouteSpec_To_route_RouteSpec(in *RouteSpec, out *route.RouteSpec, s conversion.Scope) error {
	return autoConvert_v1_RouteSpec_To_route_RouteSpec(in, out, s)
}

func autoConvert_route_RouteSpec_To_v1_RouteSpec(in *route.RouteSpec, out *RouteSpec, s conversion.Scope) error {
	out.Host = in.Host
	out.Path = in.Path
	if err := Convert_route_RouteTargetReference_To_v1_RouteTargetReference(&in.To, &out.To, s); err != nil {
		return err
	}
	out.AlternateBackends = *(*[]RouteTargetReference)(unsafe.Pointer(&in.AlternateBackends))
	out.Port = (*RoutePort)(unsafe.Pointer(in.Port))
	out.TLS = (*TLSConfig)(unsafe.Pointer(in.TLS))
	out.WildcardPolicy = WildcardPolicyType(in.WildcardPolicy)
	return nil
}

func Convert_route_RouteSpec_To_v1_RouteSpec(in *route.RouteSpec, out *RouteSpec, s conversion.Scope) error {
	return autoConvert_route_RouteSpec_To_v1_RouteSpec(in, out, s)
}

func autoConvert_v1_RouteStatus_To_route_RouteStatus(in *RouteStatus, out *route.RouteStatus, s conversion.Scope) error {
	out.Ingress = *(*[]route.RouteIngress)(unsafe.Pointer(&in.Ingress))
	return nil
}

func Convert_v1_RouteStatus_To_route_RouteStatus(in *RouteStatus, out *route.RouteStatus, s conversion.Scope) error {
	return autoConvert_v1_RouteStatus_To_route_RouteStatus(in, out, s)
}

func autoConvert_route_RouteStatus_To_v1_RouteStatus(in *route.RouteStatus, out *RouteStatus, s conversion.Scope) error {
	if in.Ingress == nil {
		out.Ingress = make([]RouteIngress, 0)
	} else {
		out.Ingress = *(*[]RouteIngress)(unsafe.Pointer(&in.Ingress))
	}
	return nil
}

func Convert_route_RouteStatus_To_v1_RouteStatus(in *route.RouteStatus, out *RouteStatus, s conversion.Scope) error {
	return autoConvert_route_RouteStatus_To_v1_RouteStatus(in, out, s)
}

func autoConvert_v1_RouteTargetReference_To_route_RouteTargetReference(in *RouteTargetReference, out *route.RouteTargetReference, s conversion.Scope) error {
	out.Kind = in.Kind
	out.Name = in.Name
	out.Weight = (*int32)(unsafe.Pointer(in.Weight))
	return nil
}

func Convert_v1_RouteTargetReference_To_route_RouteTargetReference(in *RouteTargetReference, out *route.RouteTargetReference, s conversion.Scope) error {
	return autoConvert_v1_RouteTargetReference_To_route_RouteTargetReference(in, out, s)
}

func autoConvert_route_RouteTargetReference_To_v1_RouteTargetReference(in *route.RouteTargetReference, out *RouteTargetReference, s conversion.Scope) error {
	out.Kind = in.Kind
	out.Name = in.Name
	out.Weight = (*int32)(unsafe.Pointer(in.Weight))
	return nil
}

func Convert_route_RouteTargetReference_To_v1_RouteTargetReference(in *route.RouteTargetReference, out *RouteTargetReference, s conversion.Scope) error {
	return autoConvert_route_RouteTargetReference_To_v1_RouteTargetReference(in, out, s)
}

func autoConvert_v1_RouterShard_To_route_RouterShard(in *RouterShard, out *route.RouterShard, s conversion.Scope) error {
	out.ShardName = in.ShardName
	out.DNSSuffix = in.DNSSuffix
	return nil
}

func Convert_v1_RouterShard_To_route_RouterShard(in *RouterShard, out *route.RouterShard, s conversion.Scope) error {
	return autoConvert_v1_RouterShard_To_route_RouterShard(in, out, s)
}

func autoConvert_route_RouterShard_To_v1_RouterShard(in *route.RouterShard, out *RouterShard, s conversion.Scope) error {
	out.ShardName = in.ShardName
	out.DNSSuffix = in.DNSSuffix
	return nil
}

func Convert_route_RouterShard_To_v1_RouterShard(in *route.RouterShard, out *RouterShard, s conversion.Scope) error {
	return autoConvert_route_RouterShard_To_v1_RouterShard(in, out, s)
}

func autoConvert_v1_TLSConfig_To_route_TLSConfig(in *TLSConfig, out *route.TLSConfig, s conversion.Scope) error {
	out.Termination = route.TLSTerminationType(in.Termination)
	out.Certificate = in.Certificate
	out.Key = in.Key
	out.CACertificate = in.CACertificate
	out.DestinationCACertificate = in.DestinationCACertificate
	out.InsecureEdgeTerminationPolicy = route.InsecureEdgeTerminationPolicyType(in.InsecureEdgeTerminationPolicy)
	return nil
}

func Convert_v1_TLSConfig_To_route_TLSConfig(in *TLSConfig, out *route.TLSConfig, s conversion.Scope) error {
	return autoConvert_v1_TLSConfig_To_route_TLSConfig(in, out, s)
}

func autoConvert_route_TLSConfig_To_v1_TLSConfig(in *route.TLSConfig, out *TLSConfig, s conversion.Scope) error {
	out.Termination = TLSTerminationType(in.Termination)
	out.Certificate = in.Certificate
	out.Key = in.Key
	out.CACertificate = in.CACertificate
	out.DestinationCACertificate = in.DestinationCACertificate
	out.InsecureEdgeTerminationPolicy = InsecureEdgeTerminationPolicyType(in.InsecureEdgeTerminationPolicy)
	return nil
}

func Convert_route_TLSConfig_To_v1_TLSConfig(in *route.TLSConfig, out *TLSConfig, s conversion.Scope) error {
	return autoConvert_route_TLSConfig_To_v1_TLSConfig(in, out, s)
}