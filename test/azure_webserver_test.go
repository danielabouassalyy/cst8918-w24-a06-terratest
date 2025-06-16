package test

import (
    "testing"
    "time"

    "github.com/gruntwork-io/terratest/modules/azure"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

const subscriptionID = "805ef5cd-dba2-4928-a666-50cf5f429a0d"
const labelPrefix   = "abou0344"

// helper to build our common Options
func terraformOpts() *terraform.Options {
    return &terraform.Options{
        TerraformDir: "../",
        Vars: map[string]interface{}{
            "labelPrefix": labelPrefix,
        },
        RetryableTerraformErrors: map[string]string{
            // Retry on any NIC‐update error
            "waiting for update of Network Interface": ".*",
        },
        MaxRetries:         10,
        TimeBetweenRetries: 15 * time.Second,
    }
}

func TestAzureLinuxVMCreation(t *testing.T) {
    opts := terraformOpts()
    defer terraform.Destroy(t, opts)
    terraform.InitAndApply(t, opts)

    vmName := terraform.Output(t, opts, "vm_name")
    rgName := terraform.Output(t, opts, "resource_group_name")

    assert.True(t, azure.VirtualMachineExists(t, vmName, rgName, subscriptionID))
}

func TestNICAttachedToVM(t *testing.T) {
    opts := terraformOpts()
    defer terraform.Destroy(t, opts)
    terraform.InitAndApply(t, opts)

    vmName := terraform.Output(t, opts, "vm_name")
    rgName := terraform.Output(t, opts, "resource_group_name")
    nicName := terraform.Output(t, opts, "nic_name")

    nic, err := azure.GetNetworkInterfaceE(rgName, nicName, subscriptionID)
    assert.NoError(t, err)
    assert.NotNil(t, nic.VirtualMachine, "NIC should be attached to a VM")
    assert.Contains(t, *nic.VirtualMachine.ID, vmName)
}

func TestVMOSVersion(t *testing.T) {
    opts := terraformOpts()
    defer terraform.Destroy(t, opts)
    terraform.InitAndApply(t, opts)

    vmName := terraform.Output(t, opts, "vm_name")
    rgName := terraform.Output(t, opts, "resource_group_name")

    vm, err := azure.GetVirtualMachineE(rgName, vmName, subscriptionID)
    assert.NoError(t, err)

    // Make sure this matches your Terraform source_image_reference.sku
    expectedSKU := "22_04-lts-gen2"
    sku := *vm.StorageProfile.ImageReference.Sku
    assert.Equal(t, expectedSKU, sku)
}
