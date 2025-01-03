// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;
using Pulumi;

namespace PulumiSdk.Neon.Resource
{
    [NeonResourceType("neon:resource:Project")]
    public partial class Project : global::Pulumi.CustomResource
    {
        /// <summary>
        /// URI to connect to the default database using the default endpoint.
        /// </summary>
        [Output("connection_uri")]
        public Output<string> Connection_uri { get; private set; } = null!;

        /// <summary>
        /// URI to connect to the default database using the default endpoint in the pooler mode.
        /// </summary>
        [Output("connection_uri_pooler")]
        public Output<string> Connection_uri_pooler { get; private set; } = null!;

        /// <summary>
        /// Neon default branch's name.
        /// </summary>
        [Output("default_branch_name")]
        public Output<string> Default_branch_name { get; private set; } = null!;

        /// <summary>
        /// Neon default database's name.
        /// </summary>
        [Output("default_database_name")]
        public Output<string> Default_database_name { get; private set; } = null!;

        /// <summary>
        /// The default endpoint's host.
        /// </summary>
        [Output("default_endpoint_host")]
        public Output<string> Default_endpoint_host { get; private set; } = null!;

        /// <summary>
        /// The default endpoint's host with the pooler mode active.
        /// </summary>
        [Output("default_endpoint_host_pooler")]
        public Output<string> Default_endpoint_host_pooler { get; private set; } = null!;

        /// <summary>
        /// Neon default role's name.
        /// </summary>
        [Output("default_role_name")]
        public Output<string> Default_role_name { get; private set; } = null!;

        /// <summary>
        /// Neon default role's password.
        /// </summary>
        [Output("default_role_password")]
        public Output<string> Default_role_password { get; private set; } = null!;

        /// <summary>
        /// Project ID.
        /// </summary>
        [Output("identifier")]
        public Output<string> Identifier { get; private set; } = null!;

        /// <summary>
        /// Neon project name.
        /// </summary>
        [Output("name")]
        public Output<string?> Name { get; private set; } = null!;

        /// <summary>
        /// Neon Org ID.
        /// </summary>
        [Output("org_id")]
        public Output<string?> Org_id { get; private set; } = null!;


        /// <summary>
        /// Create a Project resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Project(string name, ProjectArgs? args = null, CustomResourceOptions? options = null)
            : base("neon:resource:Project", name, args ?? new ProjectArgs(), MakeResourceOptions(options, ""))
        {
        }

        private Project(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("neon:resource:Project", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
                PluginDownloadURL = "https://github.com/kislerdm/pulumi-neon/releases/download/v${VERSION}",
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing Project resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static Project Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new Project(name, id, options);
        }
    }

    public sealed class ProjectArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Neon project name.
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        /// <summary>
        /// Neon Org ID.
        /// </summary>
        [Input("org_id")]
        public Input<string>? Org_id { get; set; }

        public ProjectArgs()
        {
        }
        public static new ProjectArgs Empty => new ProjectArgs();
    }
}
