import { Activity } from '../services/activity.service';

export interface PeakSummaries {
  peak_id: number;
  peak_name: string;
  activities: Activity[];
}
