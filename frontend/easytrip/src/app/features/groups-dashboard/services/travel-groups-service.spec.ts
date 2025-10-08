import { TestBed } from '@angular/core/testing';

import { TravelGroupsService } from './travel-groups-service';

describe('TravelGroupsService', () => {
  let service: TravelGroupsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TravelGroupsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
