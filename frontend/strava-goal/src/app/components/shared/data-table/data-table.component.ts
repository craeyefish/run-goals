import { Component, Input, Output, EventEmitter, computed, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

export type ColumnType = 'text' | 'number' | 'date' | 'link' | 'badge' | 'progress' | 'action' | 'custom';

export interface TableColumn<T = any> {
    header: string;
    field?: keyof T | string; // Field name in data object
    type: ColumnType;
    sortable?: boolean;
    width?: string; // e.g., '120px', '20%'
    align?: 'left' | 'center' | 'right';

    // Formatter for text/number columns
    formatter?: (value: any, row: T) => string;

    // For link columns
    linkFn?: (row: T) => string; // Returns URL
    linkExternal?: boolean; // Open in new tab?

    // For badge columns
    badgeClass?: (value: any, row: T) => string; // Returns CSS class

    // For progress columns
    progressValue?: (row: T) => number; // Returns 0-100
    progressLabel?: (row: T) => string;

    // For custom rendering
    customRender?: (row: T) => string; // Returns HTML string
}

export interface TableConfig {
    emptyMessage?: string;
    rowClickable?: boolean;
    striped?: boolean;
    hoverable?: boolean;
}

@Component({
    selector: 'app-data-table',
    standalone: true,
    imports: [CommonModule, RouterModule],
    templateUrl: './data-table.component.html',
    styleUrls: ['./data-table.component.scss']
})
export class DataTableComponent<T = any> {
    @Input() columns: TableColumn<T>[] = [];
    @Input() data: T[] = [];
    @Input() config: TableConfig = {
        emptyMessage: 'No data available',
        rowClickable: false,
        striped: true,
        hoverable: true
    };

    @Output() rowClick = new EventEmitter<T>();

    // Sorting state
    sortColumn = signal<string | null>(null);
    sortDirection = signal<'asc' | 'desc'>('asc');

    // Sorted data
    sortedData = computed(() => {
        const column = this.sortColumn();
        const direction = this.sortDirection();

        if (!column) {
            return this.data;
        }

        const sorted = [...this.data].sort((a, b) => {
            const aVal = this.getFieldValue(a, column);
            const bVal = this.getFieldValue(b, column);

            if (aVal === null || aVal === undefined) return 1;
            if (bVal === null || bVal === undefined) return -1;

            let comparison = 0;
            if (typeof aVal === 'string' && typeof bVal === 'string') {
                comparison = aVal.localeCompare(bVal);
            } else if (typeof aVal === 'number' && typeof bVal === 'number') {
                comparison = aVal - bVal;
            } else {
                comparison = String(aVal).localeCompare(String(bVal));
            }

            return direction === 'asc' ? comparison : -comparison;
        });

        return sorted;
    });

    onHeaderClick(column: TableColumn<T>) {
        if (!column.sortable || !column.field) return;

        const currentSort = this.sortColumn();
        const field = String(column.field);

        if (currentSort === field) {
            // Toggle direction
            this.sortDirection.set(this.sortDirection() === 'asc' ? 'desc' : 'asc');
        } else {
            // New column
            this.sortColumn.set(field);
            this.sortDirection.set('asc');
        }
    }

    onRowClick(row: T) {
        if (this.config.rowClickable) {
            this.rowClick.emit(row);
        }
    }

    getFieldValue(row: T, field: string): any {
        const keys = field.split('.');
        let value: any = row;
        for (const key of keys) {
            value = value?.[key];
        }
        return value;
    }

    getCellContent(row: T, column: TableColumn<T>): string {
        const value = column.field ? this.getFieldValue(row, String(column.field)) : null;

        switch (column.type) {
            case 'text':
                return column.formatter ? column.formatter(value, row) : String(value ?? '');

            case 'number':
                if (column.formatter) {
                    return column.formatter(value, row);
                }
                return typeof value === 'number' ? value.toLocaleString() : String(value ?? '');

            case 'date':
                if (column.formatter) {
                    return column.formatter(value, row);
                }
                if (!value) return '';
                const date = new Date(value);
                return date.toLocaleDateString();

            case 'custom':
                return column.customRender ? column.customRender(row) : '';

            default:
                return String(value ?? '');
        }
    }

    getLink(row: T, column: TableColumn<T>): string {
        return column.linkFn ? column.linkFn(row) : '#';
    }

    getBadgeClass(row: T, column: TableColumn<T>): string {
        const value = column.field ? this.getFieldValue(row, String(column.field)) : null;
        return column.badgeClass ? column.badgeClass(value, row) : 'badge-default';
    }

    getProgressValue(row: T, column: TableColumn<T>): number {
        return column.progressValue ? column.progressValue(row) : 0;
    }

    getProgressLabel(row: T, column: TableColumn<T>): string {
        return column.progressLabel ? column.progressLabel(row) : '';
    }

    getSortIcon(column: TableColumn<T>): string {
        if (!column.sortable || !column.field) return '';

        const currentSort = this.sortColumn();
        const field = String(column.field);

        if (currentSort !== field) {
            return '↕️';
        }

        return this.sortDirection() === 'asc' ? '↑' : '↓';
    }
}
