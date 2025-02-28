export interface SummitedActivity {
  user_id: number;
  user_name: string;
  activity_id: number;
  summited_at: string; // or Date if you parse it
}

export interface PeakSummaries {
  peak_id: number;
  peak_name: string;
  summits: SummitedActivity[];
}
