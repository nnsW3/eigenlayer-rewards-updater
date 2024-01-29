import {
  ActivationDelaySet as ActivationDelaySetEvent,
  ClaimerSet as ClaimerSetEvent,
  CommissionSet as CommissionSetEvent,
  PaymentClaimed as PaymentClaimedEvent,
  PaymentUpdaterSet as PaymentUpdaterSetEvent,
  RootSubmitted as RootSubmittedEvent
} from "../generated/ClaimingManager/ClaimingManager"
import {
  ActivationDelaySet,
  ClaimerSet,
  CommissionSet,
  PaymentClaimed,
  PaymentUpdaterSet,
  RootSubmitted
} from "../generated/schema"

export function handleActivationDelaySet(event: ActivationDelaySetEvent): void {
  let entity = new ActivationDelaySet(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.oldActivationDelay = event.params.oldActivationDelay
  entity.newActivationDelay = event.params.newActivationDelay

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleClaimerSet(event: ClaimerSetEvent): void {
  let entity = new ClaimerSet(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.account = event.params.account
  entity.claimer = event.params.claimer

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleCommissionSet(event: CommissionSetEvent): void {
  let entity = new CommissionSet(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operator = event.params.operator
  entity.avs = event.params.avs
  entity.commissionBips = event.params.commissionBips

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handlePaymentClaimed(event: PaymentClaimedEvent): void {
  let entity = new PaymentClaimed(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.token = event.params.token
  entity.claimer = event.params.claimer
  entity.amount = event.params.amount

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handlePaymentUpdaterSet(event: PaymentUpdaterSetEvent): void {
  let entity = new PaymentUpdaterSet(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.oldPaymentUpdater = event.params.oldPaymentUpdater
  entity.newPaymentUpdater = event.params.newPaymentUpdater

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleRootSubmitted(event: RootSubmittedEvent): void {
  let entity = new RootSubmitted(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.root = event.params.root
  entity.paymentsCalculatedUntilTimestamp =
    event.params.paymentsCalculatedUntilTimestamp
  entity.activatedAfter = event.params.activatedAfter

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}
