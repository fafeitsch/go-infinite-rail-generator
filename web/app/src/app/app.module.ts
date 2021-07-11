import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TrackTileViewerModule } from './track-tile-viewer/track-tile-viewer.module';
import { HttpClientModule } from '@angular/common/http';

@NgModule({
  declarations: [AppComponent],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    TrackTileViewerModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
