import { Controller } from "@hotwired/stimulus";
import type { SearchConfirmEvent } from "./search_bar_controller";
import type { Product } from "../types/product";
import { assert } from "../assert";

export class ProductRowController extends Controller {
    static values = {
        isLast: Boolean
    };

    declare isLastValue: boolean;

    static targets = [
        "amountInput",
        "calories",
        "proteins",
        "fats",
        "carbs"
    ];

    declare readonly amountInputTarget: HTMLInputElement;
    declare readonly caloriesTarget: HTMLElement;
    declare readonly proteinsTarget: HTMLElement;
    declare readonly fatsTarget: HTMLElement;
    declare readonly carbsTarget: HTMLElement;

    declare product: Product | null;

    connect() {
        this.product = null;
    }

    selectProduct(e: SearchConfirmEvent) {
        this.product = e.detail as Product;

        this.updateMacros();

        if (this.isLastValue) {
            this.isLastValue = false;
            this.appendNextRow();
        }

        this.amountInputTarget.focus();
    }

    appendNextRow() {
        const template = document.getElementById("product-row-template") as HTMLTemplateElement | null;
        assert(!!template, "missing product row template");

        const newRow = template.content.cloneNode(true);
        this.element.parentElement?.appendChild(newRow);
    }

    updateMacros(): void {
        if (!this.product) { return; }

        const ratio = parseInt(this.amountInputTarget.value) / 100.0;
        const ratioOf = (n: number): string => (ratio * n).toFixed(2);

        this.caloriesTarget.innerText = ratioOf(this.product.calories);
        this.proteinsTarget.innerText = ratioOf(this.product.proteins);
        this.fatsTarget.innerText = ratioOf(this.product.fats);
        this.carbsTarget.innerText = ratioOf(this.product.carbs);
    }

    destroy(): void {
        if (this.isLastValue) {
            return;
        }

        this.element.remove();
    }
};
