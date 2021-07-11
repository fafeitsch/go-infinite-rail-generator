import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TrackTileViewerComponent } from './track-tile-viewer.component';

describe('TrackTileViewerComponent', () => {
  let component: TrackTileViewerComponent;
  let fixture: ComponentFixture<TrackTileViewerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TrackTileViewerComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TrackTileViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
