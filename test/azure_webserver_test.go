package test

import (
  "os"
  "testing"

  "github.com/gruntwork-io/terratest/modules/terraform"
  "github.com/gruntwork-io/terratest/modules/azure"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)

func TestAzureLinuxVMCreation(t *testing.T) {
  t.Parallel()
  opts := &terraform.Options{
    TerraformDir: "../",
    Vars: map[string]interface{}{
      "subscriptionId": os.Getenv("AZ_SUB_ID"),
      "labelPrefix":    "abou0344",
    },
  }
  defer terraform.Destroy(t, opts)
  terraform.InitAndApply(t, opts)

  rg     := terraform.Output(t, opts, "resource_group_name")
  vmName := terraform.Output(t, opts, "vm_name")

  assert.NotEmpty(t, rg)
  assert.Contains(t, vmName, "abou0344")
}

func TestNICAttachedToVM(t *testing.T) {
  t.Parallel()
  opts := &terraform.Options{
    TerraformDir: "../",
    Vars: map[string]interface{}{
      "subscriptionId": os.Getenv("AZ_SUB_ID"),
      "labelPrefix":    "abou0344",
    },
  }
  defer terraform.Destroy(t, opts)
  terraform.InitAndApply(t, opts)

  subID   := os.Getenv("AZ_SUB_ID")
  rg      := terraform.Output(t, opts, "resource_group_name")
  nicName := terraform.Output(t, opts, "nic_name")
  vmName  := terraform.Output(t, opts, "vm_name")

  nic, err := azure.GetNetworkInterfaceE(subID, rg, nicName)
  require.NoError(t, err)
  require.NotNil(t, nic.VirtualMachine)
  assert.Contains(t, *nic.VirtualMachine.ID, vmName)
}

func TestVMOSVersion(t *testing.T) {
  t.Parallel()
  opts := &terraform.Options{
    TerraformDir: "../",
    Vars: map[string]interface{}{
      "subscriptionId": os.Getenv("AZ_SUB_ID"),
      "labelPrefix":    "abou0344",
    },
  }
  defer terraform.Destroy(t, opts)
  terraform.InitAndApply(t, opts)

  subID  := os.Getenv("AZ_SUB_ID")
  rg     := terraform.Output(t, opts, "resource_group_name")
  vmName := terraform.Output(t, opts, "vm_name")

  vm, err := azure.GetVirtualMachineE(subID, rg, vmName)
  require.NoError(t, err)
  require.NotNil(t, vm.StorageProfile)
  require.NotNil(t, vm.StorageProfile.ImageReference)
  assert.Contains(t, *vm.StorageProfile.ImageReference.Sku, "22_04")
}
