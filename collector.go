package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"fmt"
	"log"
	"math/big"
	"time"
	"strconv"

	rhp2 "go.sia.tech/core/rhp/v2"
	"go.sia.tech/core/types"
	"go.sia.tech/hostd/api"
	"go.sia.tech/hostd/host/contracts"

)

var (
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

	
	hostdRevenuePotentialActualMonth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_actual_month", Help: "Potential revenue remaining for current month"})
	hostdRevenuePotentialNextMonth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_next_month", Help: "Potential revenue for next month"})
	hostdRevenuePotentialNext2Month = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_next_2_month", Help: "Potential revenue for next 2 month"})

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



//REVENUE FOR CURRENT MONTH
	//GET CURRENT HEIGHT
	consensusState, _ := client.Consensus()
	blockHeight := float64(consensusState.ChainIndex.Height)

	//GET REMAINING BLOCKS FOR THE CURRENT MONTH
    t := time.Now()
    year, month, _ := t.Date()
    nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
    duration := nextMonth.Sub(t)
	roundedDuration := duration.Round(10 * time.Minute)
    remainingBlocksInMonth:=(roundedDuration.Minutes())/10

	//THERE ARE 4320 BLOCKS PER MONTH
	//FIND THE INITIAL & FINAL BLOCK OF THE CURRENT MONTH
	finalBlockOfMonth := uint64(blockHeight+remainingBlocksInMonth)
	fmt.Println("Al mes actual le quedan: ", remainingBlocksInMonth)
	fmt.Println("El bloque final del mes sera aproximadamente : ", finalBlockOfMonth)
	//FILTER FOR ACTIVE CONTRACTS EXPIRING ON CURRENT MONTH
	filter := contracts.ContractFilter{
		Statuses: []contracts.ContractStatus{
			contracts.ContractStatusActive,
		},
	//	MinExpirationHeight: (initialBlockOfNextMonth), //MINHEIGHT IS THE START OF NEXTMONTH
		MaxExpirationHeight: (finalBlockOfMonth),   //  MAXHEIGHT IS THE END OF CURRENT MONTH
	}

	//TOTAL POTENTIAL REVENUE FOR ACTIVE CONTRACTS ON CURRENT MONTH
	contratos, _, err := client.Contracts(filter)
	var RevenueActualMonth float64 =0

	for _, contrato := range contratos {
		RevenueActualMonth+=convertCurrency(contrato.Usage.StorageRevenue)
		RevenueActualMonth+=convertCurrency(contrato.Usage.EgressRevenue)
		RevenueActualMonth+=convertCurrency(contrato.Usage.IngressRevenue)
		RevenueActualMonth+=convertCurrency(contrato.Usage.RPCRevenue)

	}
	
	hostdRevenuePotentialActualMonth.Set(RevenueActualMonth)


//REVENUE FOR NEXT MONTH
	//INITIAL & FINAL BLOCK OF NEXT MONTH

	initialBlockOfNextMonth := uint64(finalBlockOfMonth+1)
	finalBlockOfNextMonth := uint64(blockHeight+remainingBlocksInMonth+4320)
	
	fmt.Println("initialBlockOfNextMonth : " + strconv.FormatUint(initialBlockOfNextMonth, 10))
	fmt.Println("finalBlockOfNextMonth : " + strconv.FormatUint(finalBlockOfNextMonth, 10))

	//FILTER FOR ACTIVE CONTRACTS EXPIRING NEXT MONTH
	filter2 := contracts.ContractFilter{
		Statuses: []contracts.ContractStatus{
			contracts.ContractStatusActive,
		},
		MaxExpirationHeight: (finalBlockOfNextMonth),   //MAXHEIGHT IS THE END OF NEXT MONTH
	}
	
	//TOTAL POTENTIAL REVENUE FOR EXPIRING CONTRACTS BETWEEN ACTUAL AND NEXT MONTH
	var RevenueNextMonth float64 =0
	contratos2, _, err := client.Contracts(filter2)

	for _, contrato2 := range contratos2 {
		RevenueNextMonth+=convertCurrency(contrato2.Usage.StorageRevenue)
		RevenueNextMonth+=convertCurrency(contrato2.Usage.EgressRevenue)
		RevenueNextMonth+=convertCurrency(contrato2.Usage.IngressRevenue)
		RevenueNextMonth+=convertCurrency(contrato2.Usage.RPCRevenue)
	}
	//TOTAL REVENUE FOR ACTUAL AND NEXT MONTH
	//SUSTRACT NEXT MONTH MINUS ACTUAL MONTH 
	//GET THE NEXT MONTH
	hostdRevenuePotentialNextMonth.Set(RevenueNextMonth-RevenueActualMonth)

//REVENUE FOR NEXT 2 MONTH
	//INITIAL & FINAL BLOCK OF NEXT MONTH

	//initialBlockOfNextMonth := uint64(finalBlockOfMonth+1)
	finalBlockOfNextNextMonth := uint64(blockHeight+remainingBlocksInMonth+8640)
	
	fmt.Println("finalBlockOfNextNextMonth : " + strconv.FormatUint(finalBlockOfNextNextMonth, 10))

	//FILTER FOR ACTIVE CONTRACTS EXPIRING NEXT MONTH
	filter3 := contracts.ContractFilter{
		Statuses: []contracts.ContractStatus{
			contracts.ContractStatusActive,
		},
		MaxExpirationHeight: (finalBlockOfNextNextMonth),   //MAXHEIGHT IS THE END OF NEXT MONTH
	}
	
	//TOTAL POTENTIAL REVENUE FOR EXPIRING CONTRACTS BETWEEN ACTUAL AND NEXT 2 MONTHS
	var RevenueNextNextMonth float64 =0
	contratos3, _, err := client.Contracts(filter3)

	for _, contrato3 := range contratos3 {
		RevenueNextNextMonth+=convertCurrency(contrato3.Usage.StorageRevenue)
		RevenueNextNextMonth+=convertCurrency(contrato3.Usage.EgressRevenue)
		RevenueNextNextMonth+=convertCurrency(contrato3.Usage.IngressRevenue)
		RevenueNextNextMonth+=convertCurrency(contrato3.Usage.RPCRevenue)
	}
	potentialNext2Month := RevenueNextNextMonth-RevenueNextMonth
	hostdRevenuePotentialNext2Month.Set(potentialNext2Month)


//CONSOLE PRINT OF REVENEU PER MONTHS
	fmt.Println("PENDING REVENUE FOR CURRENT MONTH = " +strconv.FormatFloat(RevenueActualMonth, 'f', 6, 64))
	fmt.Println("EXPECTED REVENUE FOR NEXT MONTH   = " +strconv.FormatFloat(RevenueNextMonth-RevenueActualMonth, 'f', 6, 64))
	fmt.Println("EXPECTED REVENUE FOR NEXT 2 MONTH = " +strconv.FormatFloat(potentialNext2Month, 'f', 6, 64))
	

	totalRevenueALLMonths:=convertCurrency(metrics.Revenue.Potential.Storage)+convertCurrency(metrics.Revenue.Potential.Ingress)+convertCurrency(metrics.Revenue.Potential.Egress)+convertCurrency(metrics.Revenue.Potential.RPC)
	fmt.Println("EXPECTED REVENUE TOTAL            = " +strconv.FormatFloat(totalRevenueALLMonths, 'f', 6, 64))

}

