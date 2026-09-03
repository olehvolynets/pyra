import { Controller } from "@hotwired/stimulus";
import type { ActionEvent } from "@hotwired/stimulus";
import { throttle } from "@github/mini-throttle";
// import { assert } from "../assert";

export type SearchConfirmEvent = CustomEvent<any>;

export class SearchBarController extends Controller {
    static targets = [
       "input",
       "options"
    ];
    declare inputTarget: HTMLInputElement;
    declare optionsTarget: HTMLElement;

    static values = {
        method: String,
        url: String,
        param: String
    };
    declare methodValue: string;
    declare hasMethodValue: boolean;
    declare urlValue: string;
    declare hasUrlValue: boolean;
    declare paramValue: string;
    declare hasParamValue: boolean;

    static classes = ["emptyList", "selectedOption"];
    declare readonly emptyListClass: string;
    declare readonly emptyListClasses: string[];
    declare readonly selectedOptionClass: string;
    declare readonly selectedOptionClasses: string[];

    declare selectedOptionIdx: number;
    declare optionValues: any[];

    connect() {
        if (!this.hasMethodValue) {
            this.methodValue = "GET";
        }

        if (!this.hasParamValue) {
            this.paramValue = "q";
        }

        if (!this.hasUrlValue) {
            console.error("search-bar controller requires URL value");
        }

        this.selectedOptionIdx = -1;
        this.optionValues = [];
    }

    search = throttle(async (e: InputEvent) => {
        const input = e.target as HTMLInputElement;

        if (input.value === "") {
            this.hideOptions();
            return;
        }

        const params = new URLSearchParams({ [this.paramValue]: input.value });
        const fullURI = this.urlValue + "?" + params.toString();

        const response = await fetch(fullURI, { method: this.methodValue });
        const results = await response.json();

        const suggestions = results.map((x: {id: number, label: string, value: any}, idx: number) => {
            const option = document.createElement("li");
            option.dataset.value = x.id.toString();
            option.dataset.searchBarIndexParam = idx.toString();
            option.dataset.action = [
                "click->search-bar#confirm",
                "mouseover->search-bar#select",
            ].join(" ");
            option.innerText = x.label;

            this.optionValues.push(x.value);

            return option;
        });

        if (suggestions.length > 0) {
            this.selectedOptionIdx = 0;
            suggestions[0].classList.add(...this.selectedOptionClasses);
        }

        this.optionsTarget.replaceChildren(...suggestions);
    });

    hideOptions() {
        this.optionsTarget.innerHTML = "";
        this.optionsTarget.classList.add(...this.emptyListClasses);
    }

    showOptions() {
        this.optionsTarget.classList.remove(...this.emptyListClasses);
    }

    selectNext() {
        const prevIdx = this.selectedOptionIdx;
        this.selectedOptionIdx++;

        if (this.selectedOptionIdx >= this.optionsTarget.children.length) {
            this.selectedOptionIdx = 0;
        }

        this.optionsTarget.children[prevIdx].classList.remove(...this.selectedOptionClasses);
        this.optionsTarget.children[this.selectedOptionIdx].classList.add(...this.selectedOptionClasses);
    }

    selectPrev() {
        const prevIdx = this.selectedOptionIdx;
        this.selectedOptionIdx--;

        if (this.selectedOptionIdx < 0) {
            this.selectedOptionIdx = this.optionsTarget.children.length - 1;
        }

        this.optionsTarget.children[prevIdx].classList.remove(...this.selectedOptionClasses);
        this.optionsTarget.children[this.selectedOptionIdx].classList.add(...this.selectedOptionClasses);
    }

    confirm(e: InputEvent) {
        e.preventDefault();

        const selectedOption = this.optionsTarget.children[this.selectedOptionIdx] as HTMLElement;
        const label = selectedOption.innerText;
        const value = selectedOption.dataset.value;

        if (!value) {
            console.error("search-bar: no value in dataset");
            return;
        }

        this.inputTarget.value = label;
        this.hideOptions();

        const confirmEvent: SearchConfirmEvent = new CustomEvent(
            "search-confirm",
            {
                bubbles: true,
                detail: this.optionValues[this.selectedOptionIdx],
            }
        );

        this.inputTarget.dispatchEvent(confirmEvent);
    }

    select(e: ActionEvent) {
        const prevIdx = this.selectedOptionIdx;
        const newIdx = parseInt(e.params.index);
        this.selectedOptionIdx = newIdx;

        this.optionsTarget.children[prevIdx].classList.remove(...this.selectedOptionClasses);
        this.optionsTarget.children[newIdx].classList.add(...this.selectedOptionClasses);
    }

    lose(e: PointerEvent) {
        if (this.intersects(e, this.inputTarget) ||
            this.intersects(e, this.optionsTarget)) {
            return;
        }

        this.inputTarget.blur();
        this.hideOptions();
    }

    private intersects(evt: PointerEvent, el: HTMLElement): boolean {
        const { x, y, width, height } = el.getBoundingClientRect();
        if (width === 0 || height === 0) {
            return false;
        }

        const pointerX = evt.clientX;
        const pointerY = evt.clientY;

        const withinX = (pointerX >= x && pointerX <= x + width);
        const withinY = (pointerY >= y && pointerY <= y + height);

        return withinX && withinY;
    }
};
