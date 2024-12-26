package myproject;

import com.pulumi.Pulumi;
import com.pulumi.neon.resource.*;

public class App {
    public static void main(String[] args) {
        Pulumi.run(_ -> {
           new Project("myproject", ProjectArgs.builder()
                            .name("myProjectProvisionedWithPulumiJavaSDK")
                            .build());
        });
    }
}
