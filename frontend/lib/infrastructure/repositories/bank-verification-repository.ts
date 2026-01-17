// Infrastructure Layer - API Repository Implementation
// Clean Architecture: Interface Adapters

import type { BankVerificationRepository } from "@/lib/domain/use-cases/verify-bank-account"
import type { User } from "@/lib/domain/entities/user"
import type { Bank } from "@/lib/domain/entities/bank"
import type { VerificationResult } from "@/lib/domain/entities/verification"

export class ApiBankVerificationRepository implements BankVerificationRepository {
  private baseUrl: string

  constructor(baseUrl?: string) {
    this.baseUrl = baseUrl || process.env.API_BASE_URL
  }

  async verify(user: User, bank: Bank): Promise<VerificationResult> {
    try {
      const response = await fetch(`${this.baseUrl}/api/verify`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ user, bank }),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || "Verification failed")
      }

      const data = await response.json()
      return {
        success: data.success,
        message: data.message,
        timestamp: new Date(data.timestamp),
        verificationId: data.verificationId,
      }
    } catch (error) {
      return {
        success: false,
        message: error instanceof Error ? error.message : "Verification failed",
        timestamp: new Date(),
      }
    }
  }
}
