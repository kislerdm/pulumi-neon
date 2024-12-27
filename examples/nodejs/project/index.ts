import {Project, ProjectArgs} from "@dkisler/pulumi-neon/resource/project";

new Project("myproject", {name: "myProjectProvisionedWithPulumiNodejsSDK"} as ProjectArgs);
