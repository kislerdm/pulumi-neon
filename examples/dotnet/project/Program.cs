using Pulumi;
using PulumiSdk.Neon.Resource;

return await Deployment.RunAsync(() =>
{
    new Project("myproject", new ProjectArgs
    {
        Name = "myProjectProvisionedWithPulumiDotnetSDK",
    });
});
