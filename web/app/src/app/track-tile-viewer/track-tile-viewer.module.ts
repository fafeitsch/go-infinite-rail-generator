import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TrackTileViewerComponent } from './track-tile-viewer.component';
import { DragDropModule } from '@angular/cdk/drag-drop';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';

@NgModule({
  declarations: [TrackTileViewerComponent],
  exports: [TrackTileViewerComponent],
  imports: [CommonModule, DragDropModule, LeafletModule],
})
export class TrackTileViewerModule {}
