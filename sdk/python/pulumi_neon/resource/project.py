# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from .. import _utilities

__all__ = ['ProjectArgs', 'Project']

@pulumi.input_type
class ProjectArgs:
    def __init__(__self__, *,
                 name: Optional[pulumi.Input[str]] = None,
                 org_id: Optional[pulumi.Input[str]] = None):
        """
        The set of arguments for constructing a Project resource.
        :param pulumi.Input[str] name: Neon project name.
        :param pulumi.Input[str] org_id: Neon Org ID.
        """
        if name is not None:
            pulumi.set(__self__, "name", name)
        if org_id is not None:
            pulumi.set(__self__, "org_id", org_id)

    @property
    @pulumi.getter
    def name(self) -> Optional[pulumi.Input[str]]:
        """
        Neon project name.
        """
        return pulumi.get(self, "name")

    @name.setter
    def name(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "name", value)

    @property
    @pulumi.getter
    def org_id(self) -> Optional[pulumi.Input[str]]:
        """
        Neon Org ID.
        """
        return pulumi.get(self, "org_id")

    @org_id.setter
    def org_id(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "org_id", value)


class Project(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 name: Optional[pulumi.Input[str]] = None,
                 org_id: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        """
        Create a Project resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[str] name: Neon project name.
        :param pulumi.Input[str] org_id: Neon Org ID.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: Optional[ProjectArgs] = None,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Create a Project resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param ProjectArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ProjectArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 name: Optional[pulumi.Input[str]] = None,
                 org_id: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ProjectArgs.__new__(ProjectArgs)

            __props__.__dict__["name"] = name
            __props__.__dict__["org_id"] = org_id
            __props__.__dict__["connection_uri"] = None
            __props__.__dict__["connection_uri_pooler"] = None
            __props__.__dict__["default_branch_name"] = None
            __props__.__dict__["default_database_name"] = None
            __props__.__dict__["default_endpoint_host"] = None
            __props__.__dict__["default_endpoint_host_pooler"] = None
            __props__.__dict__["default_role_name"] = None
            __props__.__dict__["default_role_password"] = None
            __props__.__dict__["identifier"] = None
        super(Project, __self__).__init__(
            'neon:resource:Project',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'Project':
        """
        Get an existing Project resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = ProjectArgs.__new__(ProjectArgs)

        __props__.__dict__["connection_uri"] = None
        __props__.__dict__["connection_uri_pooler"] = None
        __props__.__dict__["default_branch_name"] = None
        __props__.__dict__["default_database_name"] = None
        __props__.__dict__["default_endpoint_host"] = None
        __props__.__dict__["default_endpoint_host_pooler"] = None
        __props__.__dict__["default_role_name"] = None
        __props__.__dict__["default_role_password"] = None
        __props__.__dict__["identifier"] = None
        __props__.__dict__["name"] = None
        __props__.__dict__["org_id"] = None
        return Project(resource_name, opts=opts, __props__=__props__)

    @property
    @pulumi.getter
    def connection_uri(self) -> pulumi.Output[str]:
        """
        URI to connect to the default database using the default endpoint.
        """
        return pulumi.get(self, "connection_uri")

    @property
    @pulumi.getter
    def connection_uri_pooler(self) -> pulumi.Output[str]:
        """
        URI to connect to the default database using the default endpoint in the pooler mode.
        """
        return pulumi.get(self, "connection_uri_pooler")

    @property
    @pulumi.getter
    def default_branch_name(self) -> pulumi.Output[str]:
        """
        Neon default branch's name.
        """
        return pulumi.get(self, "default_branch_name")

    @property
    @pulumi.getter
    def default_database_name(self) -> pulumi.Output[str]:
        """
        Neon default database's name.
        """
        return pulumi.get(self, "default_database_name")

    @property
    @pulumi.getter
    def default_endpoint_host(self) -> pulumi.Output[str]:
        """
        The default endpoint's host.
        """
        return pulumi.get(self, "default_endpoint_host")

    @property
    @pulumi.getter
    def default_endpoint_host_pooler(self) -> pulumi.Output[str]:
        """
        The default endpoint's host with the pooler mode active.
        """
        return pulumi.get(self, "default_endpoint_host_pooler")

    @property
    @pulumi.getter
    def default_role_name(self) -> pulumi.Output[str]:
        """
        Neon default role's name.
        """
        return pulumi.get(self, "default_role_name")

    @property
    @pulumi.getter
    def default_role_password(self) -> pulumi.Output[str]:
        """
        Neon default role's password.
        """
        return pulumi.get(self, "default_role_password")

    @property
    @pulumi.getter
    def identifier(self) -> pulumi.Output[str]:
        """
        Project ID.
        """
        return pulumi.get(self, "identifier")

    @property
    @pulumi.getter
    def name(self) -> pulumi.Output[Optional[str]]:
        """
        Neon project name.
        """
        return pulumi.get(self, "name")

    @property
    @pulumi.getter
    def org_id(self) -> pulumi.Output[Optional[str]]:
        """
        Neon Org ID.
        """
        return pulumi.get(self, "org_id")

