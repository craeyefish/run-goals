import { Component } from '@angular/core';
import { PeakSummitTableComponent } from 'src/app/components/peak-summit-table/peak-summit-table.component';
import { GroupsTableComponent } from '../../components/groups/groups-table/groups-table.component';

@Component({
  selector: 'app-summits',
  standalone: true,
  imports: [PeakSummitTableComponent],
  templateUrl: './summits-page.component.html',
  styleUrls: ['./summits-page.component.scss'],
})
export class SummitsPageComponent {}
