package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	//	"fmt"
	"log"
	"math/big"
	"time"

	rhp2 "go.sia.tech/core/rhp/v2"
	"go.sia.tech/core/types"
	"go.sia.tech/hostd/api"
)

var (
	/*	// Revenue Metrics
		hostStorageRevenue = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "host_storage_potential", Help: "Storage potential revenue"})
	*/
	hostdTotalStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_total_storage", Help: "Total amount of storage available on the hostd in bytes"})
	hostdUsedStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_used_storage", Help: "Total amount of storage used on the hostd in bytes"})
	hostdRemainingStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_remaining_storage", Help: "Amount of storage remaining on the host in bytes"})
	contractStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_contract_storage", Help: "Amount of contract storage on the host in bytes"})
	tempStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_temp_storage", Help: "Amount of temporary storage on the host in bytes"})

	storageReads = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_storage_reads", Help: "Amount of read operations"})
	storageWrites = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_storage_writes", Help: "Amount of write operations"})

	hostdIngress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_ingress", Help: "Total ingress bandwidth usage"})
	hostdEgress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_egress", Help: "Total egress bandwidth usage"})

	hostdLockedCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_locked_collateral", Help: "Locked collateral"})
	hostdRiskedCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_risked_collateral", Help: "Risked collateral"})

	walletConfirmedSiacoinBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_wallet_confirmed_siacoin_balance", Help: "Wallet confirmed SCP balance"})

	hostdPendingContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_pending_contract_count", Help: "Number of pending contracts"})
	hostdActiveContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_active_contract_count", Help: "Number of active contracts"})
	hostdRejectedContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_rejected_contract_count", Help: "Number of rejected contracts"})
	hostdFailedContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_failed_contract_count", Help: "Number of failed contracts"})
	hostdSuccessfulContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_successful_contract_count", Help: "Number of successful contracts"})

	hostdContractPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_contract_price", Help: "Contract price"})
	hostdIngressPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_ingress_price", Help: "Ingress price"})
	hostdEgressPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_egress_price", Help: "Egress price"})
	hostdBaseRPCPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_baserpc_price", Help: "BaseRPC price"})
	hostdSectorAccessPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_sector_access_price", Help: "SectorAccess price"})
	hostdStoragePrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_storage_price", Help: "Storage price"})
	hostdCollateralMultiplier = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_collateral_multiplier", Help: "Collateral Multiplier"})

	hostdRevenueEarnedRPC = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_earned_rpc", Help: "Revenue earned for RPC"})
	hostdRevenueEarnedStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_earned_storage", Help: "Revenue earned for storage"})
	hostdRevenueEarnedIngress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_earned_ingress", Help: "Revenue earned for ingress"})
	hostdRevenueEarnedEgress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_earned_egress", Help: "Revenue earned for egress"})
	hostdRevenueEarnedRegistryRead = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_earned_registry_read", Help: "Revenue earned for registry reads"})
	hostdRevenueEarnedRegistryWrite = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_earned_registry_write", Help: "Revenue earned for registry writes"})

	hostdRevenuePotentialRPC = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_rpc", Help: "Potential revenue for RPC"})
	hostdRevenuePotentialStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_storage", Help: "Potential revenue for storage"})
	hostdRevenuePotentialIngress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_ingress", Help: "Potential revenue for ingress"})
	hostdRevenuePotentialEgress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_egress", Help: "Potential revenue for egress"})
	hostdRevenuePotentialRegistryRead = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_registry_read", Help: "Potential revenue for registry reads"})
	hostdRevenuePotentialRegistryWrite = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_registry_write", Help: "Potential revenue for registry writes"})
)

func convertCurrency(c types.Currency) float64 {
	f, _ := new(big.Rat).SetFrac(c.Big(), types.Siacoins(1).Big()).Float64()
	return f
}

func callClient(passwd string, address string) {
	client := api.NewClient("http://"+address+"/api", passwd)
	metrics, err := client.Metrics(time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	// Storage
	hostdTotalStorage.Set(float64((metrics.Storage.TotalSectors) * rhp2.SectorSize))
	hostdUsedStorage.Set(float64((metrics.Storage.PhysicalSectors) * rhp2.SectorSize))
	contractStorage.Set(float64((metrics.Storage.ContractSectors) * rhp2.SectorSize))
	tempStorage.Set(float64((metrics.Storage.TempSectors) * rhp2.SectorSize))
	storageReads.Set(float64(metrics.Storage.Reads))
	storageWrites.Set(float64(metrics.Storage.Writes))
	hostdRemainingStorage.Set(float64((metrics.Storage.TotalSectors - metrics.Storage.PhysicalSectors) * rhp2.SectorSize))

	// Data
	hostdIngress.Set(float64(metrics.Data.RHP2.Ingress + metrics.Data.RHP3.Ingress))
	hostdEgress.Set(float64(metrics.Data.RHP2.Egress + metrics.Data.RHP3.Egress))

	// Balance
	walletConfirmedSiacoinBalance.Set(convertCurrency(metrics.Balance))

	// Contracts
	hostdLockedCollateral.Set(convertCurrency(metrics.Contracts.LockedCollateral))
	hostdRiskedCollateral.Set(convertCurrency(metrics.Contracts.RiskedCollateral))
	hostdPendingContractCount.Set(float64(metrics.Contracts.Pending))
	hostdActiveContractCount.Set(float64(metrics.Contracts.Active))
	hostdRejectedContractCount.Set(float64(metrics.Contracts.Rejected))
	hostdFailedContractCount.Set(float64(metrics.Contracts.Failed))
	hostdSuccessfulContractCount.Set(float64(metrics.Contracts.Successful))

	// Pricing
	hostdContractPrice.Set(convertCurrency(metrics.Pricing.ContractPrice))
	hostdIngressPrice.Set(convertCurrency(metrics.Pricing.IngressPrice))
	hostdEgressPrice.Set(convertCurrency(metrics.Pricing.EgressPrice))
	hostdBaseRPCPrice.Set(convertCurrency(metrics.Pricing.BaseRPCPrice))
	hostdSectorAccessPrice.Set(convertCurrency(metrics.Pricing.SectorAccessPrice))
	hostdStoragePrice.Set(convertCurrency(metrics.Pricing.StoragePrice))
	hostdCollateralMultiplier.Set(float64(metrics.Pricing.CollateralMultiplier))

	// Revenue Earned
	hostdRevenueEarnedRPC.Set(convertCurrency(metrics.Revenue.Earned.RPC))
	hostdRevenueEarnedStorage.Set(convertCurrency(metrics.Revenue.Earned.Storage))
	hostdRevenueEarnedIngress.Set(convertCurrency(metrics.Revenue.Earned.Ingress))
	hostdRevenueEarnedEgress.Set(convertCurrency(metrics.Revenue.Earned.Egress))
	hostdRevenueEarnedRegistryRead.Set(convertCurrency(metrics.Revenue.Earned.RegistryRead))
	hostdRevenueEarnedRegistryWrite.Set(convertCurrency(metrics.Revenue.Earned.RegistryWrite))

	// Revenue Potential
	hostdRevenuePotentialRPC.Set(convertCurrency(metrics.Revenue.Potential.RPC))
	hostdRevenuePotentialStorage.Set(convertCurrency(metrics.Revenue.Potential.Storage))
	hostdRevenuePotentialIngress.Set(convertCurrency(metrics.Revenue.Potential.Ingress))
	hostdRevenuePotentialEgress.Set(convertCurrency(metrics.Revenue.Potential.Egress))
	hostdRevenuePotentialRegistryRead.Set(convertCurrency(metrics.Revenue.Potential.RegistryRead))
	hostdRevenuePotentialRegistryWrite.Set(convertCurrency(metrics.Revenue.Potential.RegistryWrite))
}
