package scanner

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type WasteReport struct {
	UnattachedVolumes []string
	UnassociatedIPs   []string
	EstimatedMonthlyWaste float64
}

// Pricing constants (Average AWS Costs per month)
const (
	EBSCostPerGBPerMonth = 0.08 // $0.08 per GB for standard gp2/gp3 EBS
	EIPCostPerMonth      = 3.60 // ~$0.005/hr for unassociated Elastic IP
)

// ScanRuntimeWaste checks live AWS resources for unattached EBS volumes and unassociated Elastic IPs.
func ScanRuntimeWaste(ctx context.Context, region string) (*WasteReport, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	client := ec2.NewFromConfig(cfg)
	report := &WasteReport{}

	// 1. Scan Unattached EBS Volumes
	volInput := &ec2.DescribeVolumesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("status"),
				Values: []string{"available"},
			},
		},
	}

	volOutput, err := client.DescribeVolumes(ctx, volInput)
	if err != nil {
		return nil, fmt.Errorf("failed to describe volumes: %w", err)
	}

	for _, vol := range volOutput.Volumes {
		id := aws.ToString(vol.VolumeId)
		size := aws.ToInt32(vol.Size)
		cost := float64(size) * EBSCostPerGBPerMonth
		report.EstimatedMonthlyWaste += cost
		log.Printf("[WASTE DETECTED] Unattached EBS Volume: ID=%s, Size=%d GB, Waste=~$%.2f/mo\n", id, size, cost)
		report.UnattachedVolumes = append(report.UnattachedVolumes, id)
	}

	// 2. Scan Unassociated Elastic IPs
	eipInput := &ec2.DescribeAddressesInput{}
	eipOutput, err := client.DescribeAddresses(ctx, eipInput)
	if err != nil {
		return nil, fmt.Errorf("failed to describe Elastic IPs: %w", err)
	}

	for _, addr := range eipOutput.Addresses {
		if addr.AssociationId == nil {
			ip := aws.ToString(addr.PublicIp)
			report.EstimatedMonthlyWaste += EIPCostPerMonth
			log.Printf("[WASTE DETECTED] Unassociated Elastic IP: IP=%s, Waste=~$%.2f/mo\n", ip, EIPCostPerMonth)
			report.UnassociatedIPs = append(report.UnassociatedIPs, ip)
		}
	}

	return report, nil
}
