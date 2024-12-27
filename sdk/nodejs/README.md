# NodeJS SDK for Pulumi Neon Provider

-----

<div align="center">
    ⭐ The project needs your support! Please leave a star and become a GitHub sponsor! ⭐
</div>

-----

The SDK to provision [Neon](https://neon.tech/) Postgres projects using the [Pulumi Neon provider](https://github.com/kislerdm/pulumi-neon).

## How to use

### Prerequisites

- Pulumi v3.134.1+
- nodejs 18+
- [Neon API key](https://api-docs.neon.tech/reference/authentication#neon-api-keys) exported as env variable `NEON_API_KEY`

### Steps

1. Create a Pulumi project by running `pulumi new typescript`
2. Configure the Pulumi secret by running `pulumi config set --secret neon:api_key ${NEON_API_KEY}`.
3. Add the SDK as dependency:

```shell
npm install "@dkisler/pulumi-neon"
```

4. Edit the file `index.ts`:

```typescript
import {Project, ProjectArgs} from "@dkisler/pulumi-neon/resource/project";

new Project("myproject", {name: "myProjectProvisionedWithPulumiNodejsSDK"} as ProjectArgs);
```

5. Run `pulumi up -f`
6. Examine the Neon console: it's expected to see a new project called "myProjectProvisionedWithPulumiNodejsSDK".

## How to contribute

Please raise [GitHub issue](https://github.com/kislerdm/pulumi-neon/issues/new) in case of proposals, or found bugs.
