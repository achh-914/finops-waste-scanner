package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/achh-914/finops-waste-scanner/pkg/scanner"
)

func main() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	log.Printf("Starting FinOps Runtime Waste Scan in region: %s...\n", region)

	ctx := context.Background()
	report, err := scanner.ScanRuntimeWaste(ctx, region)
	if err != nil {
		log.Fatalf("Runtime scan execution failed: %v", err)
	}

	totalWaste := len(report.UnattachedVolumes) + len(report.UnassociatedIPs)

	fmt.Println("\n================ FINOPS FINANCIAL IMPACT REPORT ================")
	fmt.Printf("Unattached EBS Volumes:  %d\n", len(report.UnattachedVolumes))
	fmt.Printf("Unassociated Elastic IPs: %d\n", len(report.UnassociatedIPs))
	fmt.Printf("Total Idle Resources:    %d\n", totalWaste)
	fmt.Printf("ESTIMATED MONTHLY WASTE: $%.2f / month\n", report.EstimatedMonthlyWaste)
	fmt.Println("================================================================")

	// Hard Gate Enforcement: Fail pipeline if waste is detected
	if totalWaste > 0 {
		fmt.Printf("\n[FAIL] FinOps Gate Check Failed: $%.2f/month in waste detected.\n", report.EstimatedMonthlyWaste)
		fmt.Println("Please delete unattached infrastructure before merging.")
		os.Exit(1)
	}

	fmt.Println("\n[PASS] Clean Infrastructure: No active financial waste detected.")
	os.Exit(0)
}
