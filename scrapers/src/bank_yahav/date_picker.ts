import { expect, Locator, Page } from "@playwright/test"

export class DatePicker {
    private readonly locator: Locator
    private readonly inputLocator: Locator
    private readonly pmuLocator: Locator
    private readonly pmuInstanceLocator: Locator
    private readonly pmuPrevButtonLocator: Locator
    private readonly pmuNextButtonLocator: Locator
    private readonly pmuDaysLocator: Locator

    constructor(private readonly page: Page, name: string) {
        this.locator              = page.locator(`div.date-options-cell > date-picker[calendar-position="currentAccount.${name}"]`)
        this.inputLocator         = this.locator.locator(`> div.date-picker-box > input`)
        this.pmuLocator           = this.locator.locator(`> div.pickmeup`)
        this.pmuInstanceLocator   = this.pmuLocator.locator(`> div.pmu-instance`)
        this.pmuPrevButtonLocator = this.pmuInstanceLocator.locator(`> nav > div.pmu-prev`)
        this.pmuNextButtonLocator = this.pmuInstanceLocator.locator(`> nav > div.pmu-next`)
        this.pmuDaysLocator       = this.pmuInstanceLocator.locator(`> div.pmu-days`)
    }

    async open() {
        await this.inputLocator.click()
        await expect(this.pmuDaysLocator).toBeVisible()
    }

    async close() {
        await this.pmuInstanceLocator.press("Escape")
    }

    async navigateToPreviousMonth() {
        await this.pmuPrevButtonLocator.click()
        await this.page.waitForTimeout(300)
    }

    async navigateToNextMonth() {
        await this.pmuNextButtonLocator.click()
        await this.page.waitForTimeout(300)
    }

    async hasNextMonth() {
        return this.pmuNextButtonLocator.isVisible()
    }

    async hasPrevMonth() {
        return this.pmuPrevButtonLocator.isVisible()
    }

    async selectFirstDayOfMonth() {
        const firstEligibleDay = this.pmuDaysLocator.locator(`> div.pmu-button:not(.pmu-disabled):not(.pmu-not-in-month)`).first()
        await expect(firstEligibleDay).toBeVisible()
        await firstEligibleDay.click()
        await this.page.waitForTimeout(300)
    }

    async selectLastDayOfMonth() {
        const lastEligibleDay = this.pmuDaysLocator
                                    .locator(`> div.pmu-button:not(.pmu-disabled):not(.pmu-not-in-month)`)
                                    .last()
        await expect(lastEligibleDay).toBeVisible()
        await lastEligibleDay.click()
        await this.page.waitForTimeout(300)
    }
}
