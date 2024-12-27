import * as pulumi from "@pulumi/pulumi";
export declare class Project extends pulumi.CustomResource {
    /**
     * Get an existing Project resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): Project;
    /**
     * Returns true if the given object is an instance of Project.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    static isInstance(obj: any): obj is Project;
    /**
     * URI to connect to the default database using the default endpoint.
     */
    readonly connection_uri: pulumi.Output<string>;
    /**
     * URI to connect to the default database using the default endpoint in the pooler mode.
     */
    readonly connection_uri_pooler: pulumi.Output<string>;
    /**
     * Neon default branch's name.
     */
    readonly default_branch_name: pulumi.Output<string>;
    /**
     * Neon default database's name.
     */
    readonly default_database_name: pulumi.Output<string>;
    /**
     * The default endpoint's host.
     */
    readonly default_endpoint_host: pulumi.Output<string>;
    /**
     * The default endpoint's host with the pooler mode active.
     */
    readonly default_endpoint_host_pooler: pulumi.Output<string>;
    /**
     * Neon default role's name.
     */
    readonly default_role_name: pulumi.Output<string>;
    /**
     * Neon default role's password.
     */
    readonly default_role_password: pulumi.Output<string>;
    /**
     * Project ID.
     */
    readonly identifier: pulumi.Output<string>;
    /**
     * Neon project name.
     */
    readonly name: pulumi.Output<string | undefined>;
    /**
     * Neon Org ID.
     */
    readonly org_id: pulumi.Output<string | undefined>;
    /**
     * Create a Project resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args?: ProjectArgs, opts?: pulumi.CustomResourceOptions);
}
/**
 * The set of arguments for constructing a Project resource.
 */
export interface ProjectArgs {
    /**
     * Neon project name.
     */
    name?: pulumi.Input<string>;
    /**
     * Neon Org ID.
     */
    org_id?: pulumi.Input<string>;
}
