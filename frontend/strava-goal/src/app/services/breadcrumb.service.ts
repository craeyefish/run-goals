import { Injectable, signal } from "@angular/core";

@Injectable({
  providedIn: 'root',
})
export class BreadcrumbService {
  items = signal<BreadcrumbItem[]>([]);

  setItems(items: BreadcrumbItem[]) {
    this.items.set(items);
  }

  addItem(item: BreadcrumbItem) {
    this.items.update(current => [...current, item])
  }

  clear() {
    this.items.set([]);
  }
}


export interface BreadcrumbItem {
  label: string;
  routerLink?: string;
  callback?: () => void;
}
