package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

//	"fmt"
	"log"
	"math/big"
	"time"
//	"strconv"

	rhp2 "go.sia.tech/core/rhp/v2"
	"go.sia.tech/core/types"
	"go.sia.tech/hostd/v2/api"
	"go.sia.tech/hostd/host/v2/contracts"

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

	hostdRevenueDay1= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_1", Help: "Potential revenue for day 1"})
	hostdRevenueDay2= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_2", Help: "Potential revenue for day 2"})
	hostdRevenueDay3= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_3", Help: "Potential revenue for day 3"})
	hostdRevenueDay4= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_4", Help: "Potential revenue for day 4"})
	hostdRevenueDay5= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_5", Help: "Potential revenue for day 5"})
	hostdRevenueDay6= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_6", Help: "Potential revenue for day 6"})
	hostdRevenueDay7= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_7", Help: "Potential revenue for day 7"})
	hostdRevenueDay8= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_8", Help: "Potential revenue for day 8"})
	hostdRevenueDay9= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_9", Help: "Potential revenue for day 9"})
	hostdRevenueDay10= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_10", Help: "Potential revenue for day 10"})
	hostdRevenueDay11= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_11", Help: "Potential revenue for day 11"})
	hostdRevenueDay12= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_12", Help: "Potential revenue for day 12"})
	hostdRevenueDay13= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_13", Help: "Potential revenue for day 13"})
	hostdRevenueDay14= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_14", Help: "Potential revenue for day 14"})
	hostdRevenueDay15= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_15", Help: "Potential revenue for day 15"})
	hostdRevenueDay16= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_16", Help: "Potential revenue for day 16"})
	hostdRevenueDay17= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_17", Help: "Potential revenue for day 17"})
	hostdRevenueDay18= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_18", Help: "Potential revenue for day 18"})
	hostdRevenueDay19= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_19", Help: "Potential revenue for day 19"})
	hostdRevenueDay20= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_20", Help: "Potential revenue for day 20"})
	hostdRevenueDay21= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_21", Help: "Potential revenue for day 21"})
	hostdRevenueDay22= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_22", Help: "Potential revenue for day 22"})
	hostdRevenueDay23= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_23", Help: "Potential revenue for day 23"})
	hostdRevenueDay24= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_24", Help: "Potential revenue for day 24"})
	hostdRevenueDay25= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_25", Help: "Potential revenue for day 25"})
	hostdRevenueDay26= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_26", Help: "Potential revenue for day 26"})
	hostdRevenueDay27= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_27", Help: "Potential revenue for day 27"})
	hostdRevenueDay28= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_28", Help: "Potential revenue for day 28"})
	hostdRevenueDay29= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_29", Help: "Potential revenue for day 29"})
	hostdRevenueDay30= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_30", Help: "Potential revenue for day 30"})
	hostdRevenueDay31= promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_revenue_potential_day_31", Help: "Potential revenue for day 31"})
	)

func convertCurrency(c types.Currency) float64 {
	f, _ := new(big.Rat).SetFrac(c.Big(), types.Siacoins(1).Big()).Float64()
	return f
}

func calcEarningsPerDay(client *api.Client, blockHeight float64) {
	//GET REMAINING BLOCKS FOR THE CURRENT DAY 
	t := time.Now() 
	year, month, day := t.Date() 
	
	
	nextDay := time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC) 
	duration := nextDay.Sub(t) 

	roundedDuration := duration.Round(10 * time.Minute) 
	remainingBlocksInDay:=(roundedDuration.Minutes())/10 
	// Use remainingBlocksInDay instead of remainingBlocksInMonth
	finalBlockOfToday := uint64(blockHeight+remainingBlocksInDay)
	//144 BLOCKS PER DAY
	//scan the next 30 days for earnings of every day
	var revenueArray[31]float64
	var revenueDia[31]float64

//RECORREMOS DIA A DIA Y OBTENEMOS LOS CONTRATOS QUE FINALIZAN CADA DIA
	for dia := 31; dia >= 1; dia-- {
		nextDayFinalBlock := finalBlockOfToday+uint64(dia*144)
	//	fmt.Println("final Block Of DAY [" + strconv.Itoa(dia) +"] = " + strconv.FormatUint(nextDayFinalBlock, 10))

		filter := contracts.ContractFilter{
			Statuses: []contracts.ContractStatus{
				contracts.ContractStatusActive,
			},
		//	MinExpirationHeight: (initialBlockOfNextMonth), //MINHEIGHT IS THE START OF NEXTMONTH
			MaxExpirationHeight: (nextDayFinalBlock),   //  MAXHEIGHT IS THE END OF CURRENT MONTH
		}
		contratos, _, _ := client.Contracts(filter)
		
		var RevenuePerDay float64 = 0 
		
	//	RECORREMOS TODOS LOS CONTRATOS DE CADDA DIA Y SUMAMOS LAS GANANCIAS
		for _, contrato := range contratos {
			RevenuePerDay+=convertCurrency(contrato.Usage.StorageRevenue)
			RevenuePerDay+=convertCurrency(contrato.Usage.EgressRevenue)
			RevenuePerDay+=convertCurrency(contrato.Usage.IngressRevenue)
			RevenuePerDay+=convertCurrency(contrato.Usage.RPCRevenue)
		}
		
		revenueArray[dia-1] = RevenuePerDay

	}
//	fmt.Println(revenueArray)
	//OBTENEMOS LAS GANANCIAS LIMPIAS DE CADA DIA
	for dia := 31; dia >= 1; dia-- {
		if dia > 1{
    		revenueDia[dia-1]=revenueArray[dia-1]-revenueArray[dia-2]
		}else{
			revenueDia[dia-1]=revenueArray[dia-1]
		}
	}
	hostdRevenueDay1.Set(revenueDia[0])
	hostdRevenueDay2.Set(revenueDia[1])
	hostdRevenueDay3.Set(revenueDia[2])
	hostdRevenueDay4.Set(revenueDia[3])
	hostdRevenueDay5.Set(revenueDia[4])
	hostdRevenueDay6.Set(revenueDia[5])
	hostdRevenueDay7.Set(revenueDia[6])
	hostdRevenueDay8.Set(revenueDia[7])
	hostdRevenueDay9.Set(revenueDia[8])
	hostdRevenueDay10.Set(revenueDia[9])
	hostdRevenueDay11.Set(revenueDia[10])
	hostdRevenueDay12.Set(revenueDia[11])
	hostdRevenueDay13.Set(revenueDia[12])
	hostdRevenueDay14.Set(revenueDia[13])
	hostdRevenueDay15.Set(revenueDia[14])
	hostdRevenueDay16.Set(revenueDia[15])
	hostdRevenueDay17.Set(revenueDia[16])
	hostdRevenueDay18.Set(revenueDia[17])
	hostdRevenueDay19.Set(revenueDia[18])
	hostdRevenueDay20.Set(revenueDia[19])
	hostdRevenueDay21.Set(revenueDia[20])
	hostdRevenueDay22.Set(revenueDia[21])
	hostdRevenueDay23.Set(revenueDia[22])
	hostdRevenueDay24.Set(revenueDia[23])
	hostdRevenueDay25.Set(revenueDia[24])
	hostdRevenueDay26.Set(revenueDia[25])
	hostdRevenueDay27.Set(revenueDia[26])
	hostdRevenueDay28.Set(revenueDia[27])
	hostdRevenueDay29.Set(revenueDia[28])
	hostdRevenueDay30.Set(revenueDia[29])
	hostdRevenueDay31.Set(revenueDia[30])
	
//	fmt.Println(revenueDia)



}

func callClient(passwd string, address string) {
	client := api.NewClient("http://"+address+"/api", passwd)
	metrics, err := client.Metrics(time.Now())
	

	if err != nil {
		log.Fatalln(err)
	}

//METRICS
	// Storage
	hostdTotalStorage.Set(float64((metrics.Storage.TotalSectors) * rhp2.SectorSize))
	hostdUsedStorage.Set(float64((metrics.Storage.PhysicalSectors) * rhp2.SectorSize))
	contractStorage.Set(float64((metrics.Storage.ContractSectors) * rhp2.SectorSize))
	tempStorage.Set(float64((metrics.Storage.TempSectors) * rhp2.SectorSize))
	storageReads.Set(float64(metrics.Storage.Reads))
	storageWrites.Set(float64(metrics.Storage.Writes))
	hostdRemainingStorage.Set(float64((metrics.Storage.TotalSectors - metrics.Storage.PhysicalSectors) * rhp2.SectorSize))

	// Data
	hostdIngress.Set(float64(metrics.Data.RHP.Ingress))
	hostdEgress.Set(float64(metrics.Data.RHP.Egress))

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

	calcEarningsPerDay(client , blockHeight)

	for _, contrato := range contratos {
		RevenueActualMonth+=convertCurrency(contrato.Usage.StorageRevenue)
		RevenueActualMonth+=convertCurrency(contrato.Usage.EgressRevenue)
		RevenueActualMonth+=convertCurrency(contrato.Usage.IngressRevenue)
		RevenueActualMonth+=convertCurrency(contrato.Usage.RPCRevenue)
	}
	

	hostdRevenuePotentialActualMonth.Set(RevenueActualMonth)

//REVENUE FOR NEXT MONTH
	//INITIAL & FINAL BLOCK OF NEXT MONTH
//	initialBlockOfNextMonth := uint64(finalBlockOfMonth+1)
	finalBlockOfNextMonth := uint64(blockHeight+remainingBlocksInMonth+4320)
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


//	totalRevenueALLMonths:=convertCurrency(metrics.Revenue.Potential.Storage)+convertCurrency(metrics.Revenue.Potential.Ingress)+convertCurrency(metrics.Revenue.Potential.Egress)+convertCurrency(metrics.Revenue.Potential.RPC)

}

