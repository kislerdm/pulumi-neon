# C# / .Net SDK for Pulumi Neon Provider

-----

<div align="center">
    ⭐ The project needs your support! Please leave a star and become a GitHub sponsor! ⭐
</div>

-----

The SDK to provision [Neon](https://neon.tech/) Postgres projects using the [Pulumi Neon provider](https://github.com/kislerdm/pulumi-neon).

## How to use

### Prerequisites

- Pulumi v3.134.1+
- dotnet 6+.
- [Neon API key](https://api-docs.neon.tech/reference/authentication#neon-api-keys) exported as env variable `NEON_API_KEY`

### Steps

1. Create a Pulumi project by running `pulumi new csharp`
2. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.
3. Add the SDK as dependency:

```commandline
dotnet add package PulumiSdk.Neon
```

4. Edit the file `Program.cs`:

```csharp
using Pulumi;
using PulumiSdk.Neon.Resource;

return await Deployment.RunAsync(() =>
{
    new Project("myproject", new ProjectArgs
    {
        Name = "myProjectProvisionedWithPulumiDotnetSDK",
    });
});
```

5. (Optionally) Edit the project configuration to match the .Net version of your environment:

```xml

<Project Sdk="Microsoft.NET.Sdk">

    <PropertyGroup>
        <OutputType>Exe</OutputType>
        <!-- Configure your .Ner version here -->
        <TargetFramework>net8.0</TargetFramework>
        <Nullable>enable</Nullable>
    </PropertyGroup>

    <ItemGroup>
        <PackageReference Include="Pulumi" Version="3.*"/>
        <PackageReference Include="PulumiSdk.Neon" Version="0.*"/>
    </ItemGroup>

</Project>
```

6. Run `pulumi up -f`
7. Examine the Neon console: it's expected to see a new project called "myProjectProvisionedWithPulumiDotnetSDK".

## How to contribute

Please raise [GitHub issue](https://github.com/kislerdm/pulumi-neon/issues/new) in case of proposals, or found bugs.
